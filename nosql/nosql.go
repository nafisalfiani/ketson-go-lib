package nosql

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Config struct {
	MaxIdleConn uint64        `env:"MAX_IDLE_CONN"`
	MaxIdleTime time.Duration `env:"MAX_IDLE_TIME"`
	DSN         string        `env:"DSN"`
}

func Init(cfg Config) (*mongo.Client, error) {
	opts := options.Client().
		ApplyURI(cfg.DSN).
		SetServerAPIOptions(options.ServerAPI(options.ServerAPIVersion1)).
		SetMaxPoolSize(cfg.MaxIdleConn).
		SetMinPoolSize(1).
		SetMaxConnIdleTime(cfg.MaxIdleTime)

	log.Println(cfg.DSN)
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	ctx, cleanup := context.WithTimeout(context.Background(), cfg.MaxIdleTime)
	defer cleanup()

	if err := client.Ping(ctx, readpref.Nearest()); err != nil {
		return nil, err
	}

	return client, nil
}
