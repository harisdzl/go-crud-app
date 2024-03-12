package mongo

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/harisquqo/quqo-challenge-1/domain/repository/search_repository"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/config"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


type mongoRepo struct {
	p *base.Persistence
	c *gin.Context
}


func (m mongoRepo) InsertDoc(collectionName string, value interface{}) (error) {
	// Check if there is a Mongo connection 
	DBName := config.Configuration.GetString("mongoDb.dev.name")
	if m.p.DbMongo == nil {
		return errors.New("MONGO NOT FOUND")
	}

	ctx := context.TODO()

	_, err := m.p.DbMongo.Database(DBName).Collection(collectionName).InsertOne(ctx, value)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil

}

func (m mongoRepo) UpdateDoc(productId uint, collectionName string, updatedFields interface{}) (error) {
	// Check if there is a Mongo connection 
	DBName := config.Configuration.GetString("mongoDb.dev.name")
	if m.p.DbMongo == nil {
		return errors.New("MONGO NOT FOUND")
	}
	ctx := context.TODO()
	filter := bson.M{"id": productId}

	// Define the update operation
	update := bson.M{"$set": updatedFields}

	// Perform the update
	_, err := m.p.DbMongo.Database(DBName).Collection(collectionName).UpdateOne(ctx, filter, update)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return err

}

func (m mongoRepo) DeleteSingleDoc(fieldName string, collectionName string, id int64) (error) {
	// Check if there is a Mongo connection 
	DBName := config.Configuration.GetString("mongoDb.dev.name")
	if m.p.DbMongo == nil {
		return errors.New("MONGO NOT FOUND")
	}

	ctx := context.TODO()
	filter := bson.M{fieldName: id}
	_, err := m.p.DbMongo.Database(DBName).Collection(collectionName).DeleteOne(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil

}

func (m mongoRepo) DeleteMultipleDoc(fieldName string, collectionName string, id int64) (error) {
	// Check if there is a Mongo connection 
	DBName := config.Configuration.GetString("mongoDb.dev.name")
	if m.p.DbMongo == nil {
		return errors.New("MONGO NOT FOUND")
	}

	ctx := context.TODO()
	filter := bson.M{fieldName: id}
	_, err := m.p.DbMongo.Database(DBName).Collection(collectionName).DeleteMany(ctx, filter)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil

}

func (m mongoRepo) DeleteAllDoc(collectionName string, value []interface{}) (error) {    
	// Check if there is a Mongo connection 
	DBName := config.Configuration.GetString("mongoDb.dev.name")
	if m.p.DbMongo == nil {
		return errors.New("MONGO NOT FOUND")
	}

	ctx := context.TODO()
	filter := bson.D{}
	_, err := m.p.DbMongo.Database(DBName).Collection(collectionName).DeleteMany(ctx, filter)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil
}

func (m mongoRepo) InsertAllDoc(collectionName string, value []interface{}) (error) {


	// Check if there is a Mongo connection 
	DBName := config.Configuration.GetString("mongoDb.dev.name")
	if m.p.DbMongo == nil {
		return errors.New("MONGO NOT FOUND")
	}

	ctx := context.TODO()

	_, err := m.p.DbMongo.Database(DBName).Collection(collectionName).InsertMany(ctx, value)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	return nil
}

func (m mongoRepo) SearchDocByName(name string, indexName string, src interface{}) error {	
	// Check if there is a Mongo connection 
	DBName := config.Configuration.GetString("mongoDb.dev.name")
	if m.p.DbMongo == nil {
		return errors.New("MONGO NOT FOUND")
	}
	
	ctx := context.TODO()
	searchStage := bson.D{
		{"$search", bson.D{
			{"index", indexName},
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
	cursor, err := m.p.DbMongo.Database(DBName).Collection(indexName).Aggregate(ctx, mongo.Pipeline{searchStage, sortStage}, opts)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if err := cursor.All(context.TODO(), src); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}


func NewMongoRepository(p *base.Persistence, c *gin.Context) search_repository.SearchRepository {
	return &mongoRepo{p, c}
}