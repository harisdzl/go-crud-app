package db

import (
	"fmt"
	"log"

	"github.com/harisquqo/quqo-challenge-1/infrastructure/config"
	storage_go "github.com/supabase-community/storage-go"
)



func NewSupabaseDB() (*storage_go.Client, error) {
	// connection DB
	supabaseUrl := config.Configuration.GetString("supabase.dev.url")
	supabaseKey := config.Configuration.GetString("supabase.dev.key")
	supabase := storage_go.NewClient(supabaseUrl, supabaseKey, nil)



	if supabase == nil {
		log.Fatal("Failure to connect to Supabase")
	} else {
		fmt.Println("Connected to Supabase!")
	}



	return supabase, nil

}