package db

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"

	"github.com/harisquqo/quqo-challenge-1/infrastructure/config"
	"github.com/opensearch-project/opensearch-go"
)
func NewOpenSearchDB() (*opensearch.Client, error) {
	endpoint := config.Configuration.GetString("opensearch.dev.endpoint")
	username := config.Configuration.GetString("opensearch.dev.user")
	password := config.Configuration.GetString("opensearch.dev.password")


    client, err := opensearch.NewClient(opensearch.Config{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
        },
        Addresses: []string{endpoint},
        Username:  username,
        Password:  password,
    })

    if err != nil {
        log.Fatal(err)
    }

	fmt.Println("Pinged your deployment. You successfully connected to Opensearch!")


	return client, nil
}
