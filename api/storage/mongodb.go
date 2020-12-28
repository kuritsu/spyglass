package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kuritsu/spyglass/api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	// Timeout operations after N seconds
	connectTimeout = 5
)

// MongoDB provider
type MongoDB struct {
	initialized bool
	client      *mongo.Client
	context     context.Context
	cancelFunc  context.CancelFunc
}

// Init the db
func (p *MongoDB) Init() {
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

	p.client = client
	p.context = ctx
	p.cancelFunc = cancel

	if p.initialized {
		return
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("Failed to ping cluster: %v", err)
	}

	fmt.Println("Connected to MongoDB!")

	p.createIndexes()
	p.initialized = true
}

func (p *MongoDB) createIndexes() {
	log.Println("Creating collection indexes...")
	monitorIndexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
	}
	p.client.Database("spyglass").Collection("Monitors").
		Indexes().CreateMany(p.context, monitorIndexes, nil)

	targetIndexes := []mongo.IndexModel{
		{
			Keys:    bson.M{"id": 1},
			Options: options.Index().SetUnique(true),
		},
	}
	p.client.Database("spyglass").Collection("Targets").
		Indexes().CreateMany(p.context, targetIndexes, nil)
}

// Free db connection
func (p *MongoDB) Free() {
	p.cancelFunc()
	p.client.Disconnect(p.context)
}

// GetMonitorByID returns a monitor by its ID.
func (p *MongoDB) GetMonitorByID(id string) (*types.Monitor, error) {
	col := p.client.Database("spyglass").Collection("Monitors")
	res := col.FindOne(p.context, bson.M{"id": id})
	var monitor types.Monitor
	err := res.Err()
	if err == nil {
		res.Decode(&monitor)
	} else if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &monitor, err
}

// InsertMonitor into the db
func (p *MongoDB) InsertMonitor(monitor *types.Monitor) (*types.Monitor, error) {
	monitor.CreatedAt = time.Now()
	monitor.UpdatedAt = time.Now()
	_, err := p.client.Database("spyglass").Collection("Monitors").InsertOne(
		p.context, monitor)
	if err != nil {
		log.Printf("Could not create Monitor: %v", err)
		return nil, err
	}
	return monitor, nil
}

// GetTargetByID returns a target by its ID.
func (p *MongoDB) GetTargetByID(id string) (*types.Target, error) {
	col := p.client.Database("spyglass").Collection("Targets")
	res := col.FindOne(p.context, bson.M{"id": id})
	var target types.Target
	err := res.Err()
	if err == nil {
		res.Decode(&target)
	} else if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	return &target, err
}

// InsertTarget into the db.
func (p *MongoDB) InsertTarget(target *types.Target) (*types.Target, error) {
	target.CreatedAt = time.Now()
	target.UpdatedAt = time.Now()
	_, err := p.client.Database("spyglass").Collection("Targets").InsertOne(
		p.context, target)
	if err != nil {
		log.Printf("Could not create Target: %v", err)
		return nil, err
	}
	return target, nil
}
