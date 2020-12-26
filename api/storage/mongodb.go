package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kuritsu/spyglass/api/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Timeout operations after N seconds
	connectTimeout = 5
)

// MongoDB provider
type MongoDB struct {
	client     *mongo.Client
	context    context.Context
	cancelFunc context.CancelFunc
}

// Initialize the db
func (p *MongoDB) Initialize() {
	p.createConnection()
}

func (p *MongoDB) createConnection() {
	connectionString := os.Getenv("MONGODB_CONNECTIONSTRING")
	if connectionString == "" {
		log.Fatalln("ERROR: No MongoDB connection string provided. (MONGODB_CONNECTIONSTRING)")
	}

	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Printf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		log.Printf("Failed to connect to cluster: %v", err)
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping cluster: %v", err)
	}

	fmt.Println("Connected to MongoDB!")
	p.client = client
	p.context = ctx
	p.cancelFunc = cancel
}

// GetMonitorByID returns a monitor by its ID.
func (p *MongoDB) GetMonitorByID(id string) *types.Monitor {
	return nil
}

// GetTargetByID returns a target by its ID.
func (p *MongoDB) GetTargetByID(id string) *types.Target {
	return nil
}
