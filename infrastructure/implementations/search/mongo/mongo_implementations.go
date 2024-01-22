package mongo

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/harisquqo/quqo-challenge-1/domain/repository/search_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type mongoRepo struct {
	p *base.Persistence
}


func (m mongoRepo) InsertDoc(value interface{}) (error) {
	// Check if there is a Mongo connection 
	DBName := os.Getenv("DB_MONGO_NAME")
	DbCollection := os.Getenv("DB_MONGO_COLLECTION_PRODUCTS")
	if m.p.DbMongo == nil {
		return errors.New("MONGO NOT FOUND")
	}

	ctx := context.TODO()

	_, err := m.p.DbMongo.Database(DBName).Collection(DbCollection).InsertOne(ctx, value)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil

}

func (m mongoRepo) UpdateDoc(productId uint, updatedFields interface{}) (error) {
	// Check if there is a Mongo connection 
	DBName := os.Getenv("DB_MONGO_NAME")
	DbCollection := os.Getenv("DB_MONGO_COLLECTION_PRODUCTS")
	if m.p.DbMongo == nil {
		return errors.New("MONGO NOT FOUND")
	}
	ctx := context.TODO()
	filter := bson.M{"id": productId}

	// Define the update operation
	update := bson.M{"$set": updatedFields}

	// Perform the update
	_, err := m.p.DbMongo.Database(DBName).Collection(DbCollection).UpdateOne(ctx, filter, update)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return err

}

func (m mongoRepo) DeleteSingleDoc(productId int64) (error) {
	// Check if there is a Mongo connection 
	DBName := os.Getenv("DB_MONGO_NAME")
	DbCollection := os.Getenv("DB_MONGO_COLLECTION_PRODUCTS")
	if m.p.DbMongo == nil {
		return errors.New("MONGO NOT FOUND")
	}

	ctx := context.TODO()
	filter := bson.M{"id": productId}
	_, err := m.p.DbMongo.Database(DBName).Collection(DbCollection).DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil

}

func (m mongoRepo) DeleteAllDoc(value []interface{}) (error) {
	// Check if there is a Mongo connection 
	DBName := os.Getenv("DB_MONGO_NAME")
	DbCollection := os.Getenv("DB_MONGO_COLLECTION_PRODUCTS")
	if m.p.DbMongo == nil {
		return errors.New("MONGO NOT FOUND")
	}

	ctx := context.TODO()
	filter := bson.D{}
	_, err := m.p.DbMongo.Database(DBName).Collection(DbCollection).DeleteMany(ctx, filter)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil
}

func (m mongoRepo) InsertAllDoc(value []interface{}) (error) {
	// Check if there is a Mongo connection 
	DBName := os.Getenv("DB_MONGO_NAME")
	DbCollection := os.Getenv("DB_MONGO_COLLECTION_PRODUCTS")
	if m.p.DbMongo == nil {
		return errors.New("MONGO NOT FOUND")
	}

	ctx := context.TODO()

	_, err := m.p.DbMongo.Database(DBName).Collection(DbCollection).InsertMany(ctx, value)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil
}

func (m mongoRepo) SearchDocByName(name string) (*mongo.Cursor, error) {
	// Check if there is a Mongo connection 
	DBName := os.Getenv("DB_MONGO_NAME")
	DbCollection := os.Getenv("DB_MONGO_COLLECTION_PRODUCTS")
	if m.p.DbMongo == nil {
		return nil, errors.New("MONGO NOT FOUND")
	}
	
	ctx := context.TODO()
	searchStage := bson.D{
		{"$search", bson.D{
			{"index", "default"},
			{"text", bson.D{
				{"query", name},
				{"path", bson.A{"name"}},
				{"fuzzy", bson.D{
					{"maxEdits", 2},       // Maximum edits allowed for fuzzy matching
					{"prefixLength", 2},   // Length of the common prefix required for fuzzy matching
					{"maxExpansions", 100}, // Maximum number of term expansions to be considered
				}},
			}},
		}},
	}
	
	// Additional stages based on your needs
	// Example: Sort by relevance score
	sortStage := bson.D{{"$sort", bson.D{{"score", bson.D{{"$meta", "textScore"}}}}}}

	// Use the Aggregate method to perform the search
	opts := options.Aggregate().SetMaxTime(5 * time.Second)
	cursor, err := m.p.DbMongo.Database(DBName).Collection(DbCollection).Aggregate(ctx, mongo.Pipeline{searchStage, sortStage}, opts)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return cursor, nil
}


func NewMongoRepository(p *base.Persistence) search_repository.SearchRepository {
	return &mongoRepo{p}
}