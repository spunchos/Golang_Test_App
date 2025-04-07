package mongodb

import (
	"TestApp/internal/config"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func NewClient(ctx context.Context, database string, mongoConnectionString config.Config) (db *mongo.Database, err error) {
	// host, port, username, password, , authDB
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoConnectionString.MongoDB).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("connection failed to mongodb due to error %v", err)
	}

	defer func() {
		if err = client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping failed to mongodb due to error %v", err)
	}
	fmt.Println("Successfully connected to MongoDB!")
	return client.Database(database), nil
}
