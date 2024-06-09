package mongodb

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/kuritsu/spyglass/api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

const (
	// Timeout operations after N seconds
	connectTimeout = 30
)

// MongoDB provider
type MongoDB struct {
	Log        *logrus.Logger
	client     *mongo.Client
	context    context.Context
	cancelFunc context.CancelFunc
}

// Init the db
func (p *MongoDB) Init() {
	p.Log.Debug("MongoDB.Init")
	connectionString := os.Getenv("MONGODB_CONNECTIONSTRING")
	if connectionString == "" {
		p.Log.Fatal("ERROR: No MongoDB connection string provided. (MONGODB_CONNECTIONSTRING)")
	}

	// TODO: Change this to a transaction
	client, err := mongo.NewClient(options.Client().ApplyURI(connectionString))
	if err != nil {
		p.Log.Fatalf("Failed to create client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), connectTimeout*time.Second)

	err = client.Connect(ctx)
	if err != nil {
		p.Log.Fatalf("Failed to connect to cluster: %v", err)
	}

	p.client = client
	p.context = ctx
	p.cancelFunc = cancel
}

func (p *MongoDB) Seed() {
	// Force a connection to verify our connection string
	err := p.client.Ping(p.context, nil)
	if err != nil {
		p.Log.Fatalf("Failed to ping cluster: %v", err)
	}

	p.Log.Debug("Connected to MongoDB!")

	p.Log.Println("Creating collection indexes...")
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

	jobIndexes := []mongo.IndexModel{
		{
			Keys: bson.M{"label": 1, "targetId": 2},
		},
	}
	p.client.Database("spyglass").Collection("Jobs").
		Indexes().CreateMany(p.context, jobIndexes, nil)

	schedulerIndexes := []mongo.IndexModel{
		{
			Keys: bson.M{"label": 1},
		},
	}
	p.client.Database("spyglass").Collection("Schedulers").
		Indexes().CreateMany(p.context, schedulerIndexes, nil)

	roleIndex := []mongo.IndexModel{
		{
			Keys:    bson.M{"name": 1},
			Options: options.Index().SetUnique(true),
		},
	}

	roleCollection := p.client.Database("spyglass").Collection("Roles")
	roleCollection.Indexes().CreateMany(p.context, roleIndex, nil)

	adminsRole := types.Role{
		Name:        "admins",
		Description: "Administrators",
		Permissions: types.Permissions{
			Owners:    []string{"admins"},
			Readers:   []string{"admins"},
			Writers:   []string{"admins"},
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}
	_, err = roleCollection.InsertOne(p.context, adminsRole)
	if err != nil {
		p.Log.Error(err)
	} else {
		p.Log.Info("Inserted role admins.")
	}

	userIndex := []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
	}
	userCollection := p.client.Database("spyglass").Collection("Users")
	userCollection.Indexes().CreateMany(p.context, userIndex, nil)
	epwd, _ := bcrypt.GenerateFromPassword([]byte("admin"), 14)
	adminUser := types.User{
		Email:     "admin",
		FullName:  "Administrator",
		Roles:     []string{"admins"},
		PassHash:  string(epwd),
		FirstHash: string(epwd),
		Permissions: types.Permissions{
			Owners:    []string{"admins"},
			Readers:   []string{"admins"},
			Writers:   []string{"admins"},
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}
	_, err = userCollection.InsertOne(p.context, adminUser)
	if err != nil {
		p.Log.Error(err)
	} else {
		p.Log.Info("Inserted user admin@spyglass.com.")
	}

	tokenIndex := []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1, "token": 2},
			Options: options.Index(),
		},
	}
	tokenCollection := p.client.Database("spyglass").Collection("Tokens")
	tokenCollection.Indexes().CreateMany(p.context, tokenIndex, nil)
}

// Free db connection
func (p *MongoDB) Free() {
	p.Log.Debug("MongoDB.Free")
	p.cancelFunc()
	p.client.Disconnect(p.context)
}
