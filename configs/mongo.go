package configs

import (
	"github.com/kelseyhightower/envconfig"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
	"time"
)

type MongoConfig struct {
	Uri                   string `envconfig:"URI"`
	DatabaseName          string `envconfig:"DATABASE_NAME"`
	ReadPreference        string `envconfig:"READ_PREFERENCE"`
	ConnectionTimeout     int    `envconfig:"CONNECTION_TIMEOUT"`
	SocketTimeout         int    `envconfig:"SOCKET_TIMEOUT"`
	MaxConnectionIdleTime int    `envconfig:"MAX_CONNECTION_IDLE_TIME"`
	MinPoolSize           int    `envconfig:"MIN_POOL_SIZE"`
	MaxPoolSize           int    `envconfig:"MAX_POOL_SIZE"`
	HeartbeatInterval     int    `envconfig:"HEARTBEAT_INTERVAL"`
}

func NewMongoConfig(e EnvFileRead) (MongoConfig, error) {
	var cfg MongoConfig
	if err := envconfig.Process("MONGO", &cfg); err != nil {
		return cfg, err
	}

	return cfg, nil
}

func NewMongoClient(cfg MongoConfig) (*mongo.Client, error) {
	connectTimeout := getDurationFromMilliseconds(cfg.ConnectionTimeout)
	maxConnectionIdleTime := getDurationFromMilliseconds(cfg.MaxConnectionIdleTime)
	minPoolSize := uint64(cfg.MinPoolSize)
	maxPoolSize := uint64(cfg.MaxPoolSize)
	heartbeatInterval := getDurationFromMilliseconds(cfg.HeartbeatInterval)

	clientOptions := options.ClientOptions{
		ReadPreference:    getMongoReadPreference(cfg.ReadPreference),
		ConnectTimeout:    &connectTimeout,
		MaxConnIdleTime:   &maxConnectionIdleTime,
		MinPoolSize:       &minPoolSize,
		MaxPoolSize:       &maxPoolSize,
		HeartbeatInterval: &heartbeatInterval,
	}

	return mongo.Connect(clientOptions.ApplyURI(cfg.Uri))
}

func NewMongoDatabase(cfg MongoConfig, client *mongo.Client) *mongo.Database {
	return client.Database(cfg.DatabaseName)
}

func getDurationFromMilliseconds(millisecond int) time.Duration {
	return time.Duration(millisecond) * time.Millisecond
}

func getMongoReadPreference(preference string) (r *readpref.ReadPref) {

	switch preference {
	case "PRIMARY":
		r = readpref.Primary()
	case "SECONDARY":
		r = readpref.Secondary()
	default:
		r = readpref.SecondaryPreferred()
	}

	return
}
