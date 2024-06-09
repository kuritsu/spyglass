package mongodb

import (
	"time"

	"github.com/kuritsu/spyglass/api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAllMonitors returns all monitors which contains a string, paginated.
func (p *MongoDB) GetAllMonitors(pageSize int64, pageIndex int64, contains string) ([]*types.Monitor, error) {
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
	monitors := make([]*types.Monitor, 0)
	err = cursor.All(p.context, &monitors)
	if err != nil {
		return nil, err
	}
	return monitors, nil
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
	monitor.CreatedAt = time.Now().UTC()
	monitor.UpdatedAt = time.Now().UTC()
	_, err := p.client.Database("spyglass").Collection("Monitors").InsertOne(
		p.context, monitor)
	if err != nil {
		p.Log.Errorf("Could not create Monitor: %v", err)
		return nil, err
	}
	return monitor, nil
}

// UpdateMonitor with the modifyable fields.
func (p *MongoDB) UpdateMonitor(oldMonitor *types.Monitor, newMonitor *types.Monitor) (*types.Monitor, error) {
	newMonitor.CreatedAt = oldMonitor.CreatedAt
	newMonitor.Owners = oldMonitor.Owners
	newMonitor.UpdatedAt = time.Now().UTC()
	_, err := p.client.Database("spyglass").Collection("Monitors").UpdateOne(
		p.context, bson.M{"id": oldMonitor.ID},
		bson.M{"$set": newMonitor})
	if err != nil {
		p.Log.Errorf("Could not create Monitor: %v", err)
		return nil, err
	}
	return newMonitor, nil
}
