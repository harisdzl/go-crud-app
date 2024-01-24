package db

import (
	"fmt"
	"log"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
)
func NewElasticSearchDb() (*elasticsearch.Client, error) {
	DbElasticSearchCloudID := os.Getenv("ELASTIC_SEARCH_CLOUD_ID")
	DbElasticSearchAPIKey := os.Getenv("ELASTIC_SEARCH_API_KEY")


	cfg := elasticsearch.Config{
		CloudID:DbElasticSearchCloudID,
		APIKey: DbElasticSearchAPIKey,
	}
	
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	infores, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	fmt.Println(infores)
	fmt.Println("Successfully connected to ElasticDb!")

	return es, nil
}

