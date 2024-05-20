package server

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

const (
	DB         = "logs-db"
	COLLECTION = "logs-collection"
)

type MongoRepository struct {
	conn *mongo.Client
}

func NewMongoRepository(conn *mongo.Client) *MongoRepository {
	return &MongoRepository{
		conn: conn,
	}
}

func NewMongoClient(config *MongoConfig) (*mongo.Client, context.CancelFunc, error) {
	// create connection to MongoDB
	clientOptions := options.Client().ApplyURI(config.URI)
	clientOptions.SetMaxPoolSize(uint64(config.MaxPool))
	clientOptions.SetMinPoolSize(uint64(config.MinPool))
	clientOptions.SetMaxConnIdleTime(time.Duration(config.IdleTime) * time.Minute)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	client, err := mongo.Connect(ctx, clientOptions)
	return client, cancel, err
}

func (r *MongoRepository) Add(ctx context.Context, req Log) error {
	log.Println(req)
	collection := r.conn.Database(DB).Collection(COLLECTION)
	res, err := collection.InsertOne(ctx, req)
	log.Println(res)
	if err != nil {
		return err
	}
	return nil
}

func (r *MongoRepository) GetAll(ctx context.Context) ([]Log, error) {
	collection := r.conn.Database(DB).Collection(COLLECTION)
	log.Printf("collection: %+v", collection)
	// Use an empty BSON document as the filter
	filter := bson.M{}

	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		log.Println("error finding logs")
		log.Println(err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var res []Log
	if err := cursor.All(ctx, &res); err != nil {
		log.Println("error decoding logs")
		return nil, err
	}
	return res, nil
}
