package storage

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kuritsu/spyglass/api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

	// TODO: Change this to a transaction
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

// GetAllMonitors returns all monitors which contains a string, paginated.
func (p *MongoDB) GetAllMonitors(pageSize int64, pageIndex int64, contains string) ([]types.Monitor, error) {
	col := p.client.Database("spyglass").Collection("Monitors")
	filter := bson.M{}
	if contains != "" {
		contains = types.GetIDForRegex(contains)
		filter["id"] = bson.M{"$regex": contains}
	}
	skip := pageIndex * pageSize
	opts := options.FindOptions{
		Skip:  &skip,
		Sort:  bson.M{"id": 1},
		Limit: &pageSize,
	}
	cursor, err := col.Find(p.context, filter, &opts)
	if err != nil {
		return nil, err
	}
	monitors := make([]types.Monitor, 0)
	err = cursor.All(p.context, &monitors)
	if err != nil {
		return nil, err
	}
	return monitors, nil
}

// GetAllTargets returns all targets which contains a string, paginated.
func (p *MongoDB) GetAllTargets(pageSize int64, pageIndex int64, contains string) ([]*types.Target, error) {
	col := p.client.Database("spyglass").Collection("Targets")
	containsRegex := `^[\w\d\-_]+$`
	if contains != "" {
		containsRegex = types.GetIDForRegex(contains)
	}
	skip := pageIndex * pageSize
	opts := options.FindOptions{
		Skip:  &skip,
		Sort:  bson.M{"id": 1},
		Limit: &pageSize,
	}
	cursor, err := col.Find(p.context, bson.M{"id": bson.M{"$regex": containsRegex}}, &opts)
	if err != nil {
		return nil, err
	}
	targets := make([]*types.Target, 0)
	err = cursor.All(p.context, &targets)
	if err != nil {
		return nil, err
	}
	p.updateTargetListStatus(targets)
	return targets, nil
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
func (p *MongoDB) GetTargetByID(id string, includeChildren bool) (*types.Target, error) {
	col := p.client.Database("spyglass").Collection("Targets")
	id = strings.ToLower(id)
	expr := bson.M{"id": id}
	if includeChildren {
		escapedID := types.GetIDForRegex(id)
		regx := fmt.Sprintf(`^%s(\.[\w\d_\-]+){0,1}$`, escapedID)
		expr = bson.M{"id": bson.M{"$regex": regx}}
	}
	cursor, err := col.Find(p.context, expr)
	if err != nil {
		return nil, err
	}
	var targets []*types.Target
	err = cursor.All(p.context, &targets)
	if err != nil {
		return nil, err
	}
	if len(targets) == 0 {
		return nil, nil
	}
	p.updateTargetListStatus(targets)
	var parent *types.Target
	var children []types.Target
	for _, t := range targets {
		if t.ID == id {
			parent = t
			continue
		}
		children = append(children, *t)
	}
	parent.Children = children
	return parent, err
}

// InsertTarget into the db.
func (p *MongoDB) InsertTarget(target *types.Target) (*types.Target, error) {
	target.ID = strings.ToLower(target.ID)
	target.Children = []types.Target{}
	target.CreatedAt = time.Now()
	target.UpdatedAt = time.Now()
	col := p.client.Database("spyglass").Collection("Targets")
	if _, err := col.InsertOne(p.context, target); err != nil {
		log.Printf("Could not create Target: %v", err)
		return nil, err
	}
	if err := p.updateParentStatus(target, 1, target.Status); err != nil {
		log.Printf("Error updating parents: %v", err)
		return nil, err
	}
	return target, nil
}

// UpdateTargetStatus with all modified fields
func (p *MongoDB) UpdateTargetStatus(target *types.Target, targetPatch *types.TargetPatch) (*types.Target, error) {
	col := p.client.Database("spyglass").Collection("Targets")
	lockedDoc := col.FindOneAndUpdate(p.context, bson.M{"id": target.ID},
		bson.M{"$set": bson.M{"flock": bson.M{"pseudoRandom": primitive.NewObjectID()}}})
	if lockedDoc.Err() != nil {
		return nil, lockedDoc.Err()
	}
	lockedDoc.Decode(&target)
	if targetPatch.StatusDescription == "" {
		targetPatch.StatusDescription = target.StatusDescription
	}
	_, err := col.UpdateOne(p.context, bson.M{"id": target.ID},
		bson.M{"$set": bson.M{
			"status":            targetPatch.Status,
			"statusDescription": targetPatch.StatusDescription}})
	if err != nil {
		return nil, err
	}
	diff := targetPatch.Status - target.Status
	err = p.updateParentStatus(target, 0, diff)
	if err != nil {
		return nil, err
	}
	target.Status = targetPatch.Status
	target.StatusDescription = targetPatch.StatusDescription
	return target, nil
}

func (p *MongoDB) updateParentStatus(target *types.Target, childrenCount int, statusInc int) error {
	parents := strings.Split(target.ID, ".")
	if len(parents) < 2 {
		return nil
	}
	log.Println(len(parents)-1, " parents need update.")
	col := p.client.Database("spyglass").Collection("Targets")
	parentIds := []string{}
	prefix := ""
	for idx, p := range parents[0 : len(parents)-1] {
		if idx > 0 {
			prefix = parentIds[idx-1] + "."
		}
		parentIds = append(parentIds, prefix+p)
		log.Println(parentIds[idx])
	}
	updateResult, err := col.UpdateMany(p.context,
		bson.M{"id": bson.M{"$in": parentIds}},
		bson.M{"$inc": bson.M{"childrenCount": childrenCount, "statusTotal": statusInc}})
	log.Println(updateResult)
	return err
}

func (p *MongoDB) updateTargetListStatus(targets []*types.Target) {
	for _, t := range targets {
		if t.ChildrenCount > 0 {
			if t.StatusTotal == 0 {
				t.Status = 0
				continue
			}
			t.Status = t.StatusTotal * 100 / (t.ChildrenCount * 100)
		}
	}
}
