package base

import (
	"log"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/image_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/inventory_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/product_entity"
	"github.com/harisquqo/quqo-challenge-1/domain/entity/warehouse_entity"
	"github.com/harisquqo/quqo-challenge-1/infrastructure/persistence/base/db"
	"github.com/opensearch-project/opensearch-go"
	"github.com/redis/go-redis/v9"
	storage_go "github.com/supabase-community/storage-go"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// Repositories struct -- Currently only Product Repo
type Persistence struct {
	DB *gorm.DB
	DbRedis *redis.Client
	DbMongo *mongo.Client
	DbElastic *elasticsearch.Client
	DbOpensearch *opensearch.Client
	DbSupabase *storage_go.Client
}

// Function to create a new repository
func NewPersistence() (*Persistence, error) {
	database, errDatabase := db.NewDB()
	redisDb, errRedisDb := db.NewRedisDB()
	mongoDb, errMongoDb := db.NewMongoDB()
	elasticDb, errElasticDb := db.NewElasticSearchDb()
	opensearchDb, errOpensearchDb := db.NewOpenSearchDB()
	supabaseDb, errSupabaseDb := db.NewSupabaseDB()

	if errDatabase != nil {
		log.Fatal(errDatabase)
	}

	if errRedisDb != nil {
		log.Fatal(errRedisDb)
	}

	if errMongoDb != nil {
		log.Fatal(errMongoDb)
	}

	if errElasticDb != nil {
		log.Fatal(errElasticDb)
	}

	if errOpensearchDb != nil {
		log.Fatal(errOpensearchDb)
	}

	if errSupabaseDb != nil {
		log.Fatal(errSupabaseDb)
	}



	return &Persistence{
		DB: database.DB,
		DbRedis: redisDb,
		DbMongo: mongoDb,
		DbElastic: elasticDb,
		DbOpensearch: opensearchDb,
		DbSupabase: supabaseDb,
	}, nil

}

//This migrate all tables
func (s *Persistence) Automigrate() error {
	return s.DB.AutoMigrate(&product_entity.Product{}, 
		&inventory_entity.Inventory{}, 
		&warehouse_entity.Warehouse{}, 
		&image_entity.Image{})
}