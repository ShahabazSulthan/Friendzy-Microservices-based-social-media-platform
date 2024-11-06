package db

import (
	"context"
	"fmt"
	"time"

	"github.com/ShahabazSulthan/friendzy_post/pkg/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDbCollection struct {
	OneToOneChats  *mongo.Collection
	OneToManyChats *mongo.Collection
	Chatgroup      *mongo.Collection
}

func ConnectDatabaseMongo(config *config.MongoDataBase) (*MongoDbCollection, error) {
    // Context with timeout for the connection
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Server API options
    serverAPI := options.ServerAPI(options.ServerAPIVersion1)
    fmt.Println("----------connection uri:", config.MongoDbURL)
    
    // Connecting to MongoDB
    client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoDbURL).SetServerAPIOptions(serverAPI))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to MongoDB: %w", err)
    }

    // Pinging the MongoDB server to confirm connection
    err = client.Ping(ctx, readpref.Primary())
    if err != nil {
        fmt.Println("can't ping to db, err:", err)
        return nil, fmt.Errorf("ping to MongoDB failed: %w", err)
    }

    fmt.Printf("\nconnected to MongoDB, on database %s\n", config.DataBase)

    // Initializing collections
    var mongoCollections MongoDbCollection
    db := client.Database(config.DataBase)
    mongoCollections.OneToOneChats = db.Collection("OneToOneChats")
    mongoCollections.OneToManyChats = db.Collection("OneToManyChats")
    mongoCollections.Chatgroup = db.Collection("ChatGroups")

    return &mongoCollections, nil
}

