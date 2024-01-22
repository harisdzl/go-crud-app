package base

import (
	"log"

	"github.com/harisquqo/quqo-challenge-1/domain/entity"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base/db"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// Repositories struct -- Currently only Product Repo
type Persistence struct {
	DB *gorm.DB
	DbRedis *redis.Client
	DbMongo *mongo.Client
}

// Function to create a new repository
func NewPersistence() (*Persistence, error) {
	database, errDatabase := db.NewDB()
	redisDb, errRedisDb := db.NewRedisDB()
	mongoDb, errMongoDb := db.NewMongoDB()

	if errDatabase != nil {
		log.Fatal(errDatabase)
	}

	if errRedisDb != nil {
		log.Fatal(errRedisDb)
	}

	if errMongoDb != nil {
		log.Fatal(errMongoDb)
	}




	return &Persistence{
		DB: database.DB,
		DbRedis: redisDb,
		DbMongo: mongoDb,
	}, nil

}

//This migrate all tables
func (s *Persistence) Automigrate() error {
	return s.DB.AutoMigrate(&entity.Product{})
}