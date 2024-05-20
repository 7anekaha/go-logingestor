package server

const (
	user     = "nico"
	pass     = "secret"
	hostName = "localhost"
	port     = "27017"
)

type MongoConfig struct {
	URI      string
	MinPool  uint64
	MaxPool  uint64
	IdleTime uint64
}

func NewMongoConfig(uri string, minPool, maxPool, idleTime uint64) *MongoConfig {
	return &MongoConfig{
		URI:      uri,
		MinPool:  minPool,
		MaxPool:  maxPool,
		IdleTime: idleTime,
	}
}
