package elasticSearch

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/harisquqo/quqo-challenge-1/domain/entity"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/search_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
)


type elasticSearchRepo struct {
	p *base.Persistence
}


func (e elasticSearchRepo) InsertDoc(value interface{}) (error) {
	data, _ := json.Marshal(value)
	e.p.DbElastic.Index("search-challenge", bytes.NewReader(data))
	return nil
}

func (e elasticSearchRepo) UpdateDoc(productId uint, updatedFields interface{}) (error) {

	data, _ := json.Marshal(updatedFields)

	e.p.DbElastic.Update("search-challenge", fmt.Sprint(productId), bytes.NewReader(data))
	return nil
}

func (e elasticSearchRepo) DeleteSingleDoc(productId int64) (error) {
	e.p.DbElastic.Delete("search-challenge", fmt.Sprint(productId))
	return nil
}

func (e elasticSearchRepo) DeleteAllDoc(value []interface{}) (error) {
	if e.p.DbElastic == nil {
		return errors.New("ELASTICSEARCH NOT FOUND")
	}

	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{},
		},
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		fmt.Println(err)
		return err
	}

	res, err := e.p.DbElastic.DeleteByQuery([]string{"search-challenge"}, bytes.NewReader(queryJSON))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer res.Body.Close()

	return nil
}

func (e elasticSearchRepo) InsertAllDoc(value []interface{}) (error) {
	// Iterate over each document and index it
	data, _ := json.Marshal(value)

	for _, doc := range value {
		data, err := json.Marshal(doc)

		if err != nil {
			fmt.Println(err)
			continue
		}

		// Index the document
		_, err = e.p.DbElastic.Index("search-challenge", bytes.NewReader(data))
		if err != nil {
			fmt.Println(err)
			continue
		}
	}

	e.p.DbElastic.Index("search-challenge", bytes.NewReader(data))

	return nil
}

func (e elasticSearchRepo) SearchDocByName(name string, src interface{}) error {
    query := fmt.Sprintf(`{
        "query": {
            "match": {
                "name": "%s"
            }
        }
    }`, name)

    response, err := e.p.DbElastic.Search(
        e.p.DbElastic.Search.WithIndex("search-challenge"),
        e.p.DbElastic.Search.WithBody(strings.NewReader(query)),
    )
    if err != nil {
        return err
    }

    defer response.Body.Close()

    if response.IsError() {
        return errors.New(fmt.Sprintf("Elasticsearch response error: %s", response.Status()))
    }

    var result map[string]interface{}
    if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
        return err
    }

    hits, found := result["hits"].(map[string]interface{})["hits"].([]interface{})
    if !found {
        return errors.New("unexpected response structure")
    }

    var searchResults []entity.Product

    for _, hit := range hits {
        source, ok := hit.(map[string]interface{})["_source"].(map[string]interface{})
        if !ok {
            return errors.New("unexpected response structure in _source field")
        }

        productBytes, err := json.Marshal(source)
        if err != nil {
            return err
        }

        var product entity.Product
        if err := json.Unmarshal(productBytes, &product); err != nil {
            return err
        }

        searchResults = append(searchResults, product)
    }

    reflect.ValueOf(src).Elem().Set(reflect.ValueOf(searchResults))

    return nil
}




func NewElasticSearchRepository(p *base.Persistence) search_repository.SearchRepository {
	return &elasticSearchRepo{p}
}