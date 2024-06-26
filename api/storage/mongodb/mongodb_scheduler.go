package mongodb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kuritsu/spyglass/api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (p *MongoDB) GetAllSchedulersFor(label string) ([]*types.Scheduler, error) {
	col := p.client.Database("spyglass").Collection("Schedulers")
	expr := bson.M{"label": label}
	cursor, err := col.Find(p.context, expr)
	if err != nil {
		return nil, err
	}
	var schedulers []*types.Scheduler
	err = cursor.All(p.context, &schedulers)
	if err != nil {
		return nil, err
	}
	return schedulers, nil
}

func (p *MongoDB) InsertScheduler(scheduler *types.Scheduler) (*types.Scheduler, error) {
	scheduler.ID = uuid.NewString()
	scheduler.LastPing = time.Now().UTC()

	_, err := p.client.Database("spyglass").Collection("Schedulers").InsertOne(p.context, scheduler)
	if err != nil {
		p.Log.Error(err)
		return nil, fmt.Errorf("Error creating scheduler")
	}
	return scheduler, nil
}

func (p *MongoDB) UpdateScheduler(scheduler *types.Scheduler) (*types.Scheduler, error) {
	col := p.client.Database("spyglass").Collection("Schedulers")
	scheduler.LastPing = time.Now().UTC()

	res := col.FindOneAndUpdate(p.context, bson.M{"id": scheduler.ID},
		bson.D{{"$set", scheduler}})
	if res.Err() != nil {
		p.Log.Error(res.Err().Error())
		return nil, fmt.Errorf("Error updating scheduler")
	}
	return scheduler, nil
}

func (p *MongoDB) DeleteScheduler(id string) error {
	col := p.client.Database("spyglass").Collection("Schedulers")
	_, err := col.DeleteOne(p.context, bson.M{"id": id}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (p *MongoDB) GetAllInactiveSchedulers() ([]*types.Scheduler, error) {
	col := p.client.Database("spyglass").Collection("Schedulers")
	m5 := time.Now().Add(-5 * time.Minute)
	expr := bson.M{"lastping": bson.M{"$lt": primitive.NewDateTimeFromTime(m5)}}
	cursor, err := col.Find(p.context, expr)
	if err != nil {
		return nil, err
	}
	var schedulers []*types.Scheduler
	err = cursor.All(p.context, &schedulers)
	if err != nil {
		return nil, err
	}
	return schedulers, nil
}
