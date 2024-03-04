package openSearch

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/product_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/search_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)

type opensearchRepo struct {
	p *base.Persistence
	c *gin.Context
}

func NewOpenSearchRepository(p *base.Persistence, c *gin.Context) search_repository.SearchRepository {
	return &opensearchRepo{p, c}
}

func (o *opensearchRepo) InsertDoc(indexName string, value interface{}) error {
	// Convert the value to JSON format
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Create an io.Reader from the JSON string
	reader := strings.NewReader(string(jsonValue))
	product, ok := value.(*product_entity.Product)
	if !ok {
		return errors.New("value is not of type *product_entity.Product")
	}

	// Use the OpenSearch client to index the document
	res, err := o.p.DbOpensearch.Index(
		indexName,                                  // Index name
		reader,                                     // Document body
		o.p.DbOpensearch.Index.WithRefresh("true"), // Refresh
		o.p.DbOpensearch.Index.WithDocumentID(fmt.Sprint(product.ID)), // Set document ID
	)
	if err != nil {
		log.Fatalf("ERROR: %s", err)
		return err
	}

	defer res.Body.Close()

	log.Println(res)

	return nil
}
func (o *opensearchRepo) UpdateDoc(productID uint, indexName string, updatedFields interface{}) error {
	// Convert the updated fields to JSON format
	jsonValue, err := json.Marshal(updatedFields)
	if err != nil {
		return err
	}

	// Create an io.Reader from the JSON string
	reader := strings.NewReader(string(jsonValue))
	// Use the OpenSearch client to update the document
	res, err := o.p.DbOpensearch.Update(
		indexName,           // Index name
		fmt.Sprint(productID), // Document ID
		reader,              // Updated document body
		o.p.DbOpensearch.Update.WithRefresh("true"), // Refresh
	)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Check the response status or process the response as needed
	if res.IsError() {
		return fmt.Errorf("update failed: %s", res.String())
	}

	return nil
}


func (o *opensearchRepo) DeleteSingleDoc(fieldName string, indexName string, id int64) error {
	// Use the OpenSearch client to delete the document
	_, err := o.p.DbOpensearch.Delete(indexName, fmt.Sprint(id))
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (o *opensearchRepo) DeleteMultipleDoc(fieldName string, indexName string, id int64) error {
	// Implement if needed
	return nil
}

func (o *opensearchRepo) DeleteAllDoc(indexName string, value []interface{}) error {
    ctx := context.TODO()

    // Define a query that matches all documents
    query := map[string]interface{}{
        "query": map[string]interface{}{
            "match_all": map[string]interface{}{},
        },
    }

    // Convert the query to JSON format
    queryJSON, err := json.Marshal(query)
    if err != nil {
        return err
    }

    // Create an io.Reader from the JSON string
    queryReader := strings.NewReader(string(queryJSON))
    // Use the OpenSearch client to delete documents by query
    res, err := o.p.DbOpensearch.DeleteByQuery(
		[]string{indexName},                       // Index name
        queryReader,                     // Query body
        o.p.DbOpensearch.DeleteByQuery.WithContext(ctx), // Context
    )
    if err != nil {
        log.Fatalf("ERROR: %s", err)
        return err
    }

    defer res.Body.Close()

    log.Println(res)

    return nil
}

func (o *opensearchRepo) InsertAllDoc(indexName string, value []interface{}) error {
	// Convert the value to JSON format

	for _, p := range value {
		jsonValue, err := json.Marshal(p)
		if err != nil {
			return err
		}
	
		// Create an io.Reader from the JSON string
		reader := strings.NewReader(string(jsonValue))
		product, ok := p.(product_entity.Product)
		if !ok {
			return errors.New("value is not of type *product_entity.Product")
		}

			// Use the OpenSearch client to index the document
		res, err := o.p.DbOpensearch.Index(
			indexName,                                  // Index name
			reader,                                     // Document body
			o.p.DbOpensearch.Index.WithRefresh("true"), // Refresh
			o.p.DbOpensearch.Index.WithDocumentID(fmt.Sprint(product.ID)), // Set document ID
		)

		if err != nil {
			log.Fatalf("ERROR: %s", err)
			return err
		}

		defer res.Body.Close()
	}



	return nil
}

func (o *opensearchRepo) SearchDocByName(name string, indexName string, src interface{}) error {
	// Create a query string based on the provided name
	queryString := fmt.Sprintf(`{
		"query": {
		   "bool": {
			  "should": [
				 {
					"match_phrase_prefix": {
					   "name": {
						  "query": "%s"
					   }
					}
				 },
				 {
					"fuzzy": {
					   "name": {
						  "value": "%s",
						  "fuzziness": "auto"
					   }
					}
				 },
				 {
					"query_string": {
					   "query": "*%s*",
					   "fields": ["name^3"]
					}
				 },
				 {
					"multi_match": {
					   "query": "%s",
					   "fields": ["name^3", "category^2", "description^1"]
					}
				 }
			  ],
			  "minimum_should_match": 1
		   }
		}
	 }`, name, name, name, name)
	//  queryString := fmt.Sprintf(`{
	// 	"query": {
	// 	   "bool": {
	// 		  "should": [
	// 			 {
	// 				"multi_match": {
	// 				   "query": "%s",
	// 				   "fields": ["name^3", "category^2", "description^1"], 
	// 				   "fuzziness": "AUTO"
	// 				}
	// 			 }
	// 		  ],
	// 		  "minimum_should_match": 1
	// 	   }
	// 	}
	//  }`, name)
	 

	// Create an io.Reader from the query string
	reader := strings.NewReader(queryString)

	// Use the OpenSearch client to perform the search
	res, err := o.p.DbOpensearch.Search(
		o.p.DbOpensearch.Search.WithIndex(indexName),   // Index name
		o.p.DbOpensearch.Search.WithBody(reader),        // Query body
	)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	// Check the response status or process the response as needed
	if res.IsError() {
		return fmt.Errorf("search failed: %s", res.String())
	}

	// Parse the response body and populate the src interface
	var result map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
		return err
	}

	// Extract the "hits" field from the response
	hits, ok := result["hits"].(map[string]interface{})
	if !ok {
		return errors.New("missing or invalid 'hits' field in OpenSearch response")
	}

	// Extract the "hits" array from the "hits" field
	hitsArray, ok := hits["hits"].([]interface{})
	if !ok {
		return errors.New("missing or invalid 'hits' array in OpenSearch response")
	}

	// Unmarshal each hit into the provided src interface
	var decodedHits []interface{}
	for _, hit := range hitsArray {
		source, ok := hit.(map[string]interface{})["_source"]
		if !ok {
			return errors.New("missing '_source' field in hit")
		}
		decodedHits = append(decodedHits, source)
	}

	// Marshal the decoded hits back to JSON
	decodedJSON, err := json.Marshal(decodedHits)
	if err != nil {
		return err
	}

	// Unmarshal the JSON into the provided src interface
	if err := json.Unmarshal(decodedJSON, src); err != nil {
		return err
	}

	return nil
}
