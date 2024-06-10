package mongodb

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kuritsu/spyglass/api/types"
	"go.mongodb.org/mongo-driver/bson"
)

func (p *MongoDB) InsertJob(job *types.Job) (*types.Job, error) {
	job.ID = uuid.NewString()
	job.UpdatedAt = time.Now().UTC()

	_, err := p.client.Database("spyglass").Collection("Jobs").InsertOne(p.context, job)
	if err != nil {
		p.Log.Error(err)
		return nil, fmt.Errorf("Error creating job")
	}
	return job, nil
}

func (p *MongoDB) UpdateJob(job *types.Job) (*types.Job, error) {
	col := p.client.Database("spyglass").Collection("Jobs")
	job.UpdatedAt = time.Now().UTC()

	res := col.FindOneAndUpdate(p.context, bson.M{"id": job.ID},
		bson.D{{"$set", job}})
	if res.Err() != nil {
		p.Log.Error(res.Err().Error())
		return nil, fmt.Errorf("Error updating job")
	}
	return job, nil
}

func (p *MongoDB) GetAllJobsFor(label string) ([]*types.Job, error) {
	col := p.client.Database("spyglass").Collection("Jobs")
	expr := bson.M{"label": label}
	cursor, err := col.Find(p.context, expr)
	if err != nil {
		return nil, err
	}
	var jobs []*types.Job
	err = cursor.All(p.context, &jobs)
	if err != nil {
		return nil, err
	}
	return jobs, nil
}

func (p *MongoDB) DeleteJob(id string) error {
	col := p.client.Database("spyglass").Collection("Jobs")
	_, err := col.DeleteOne(p.context, bson.M{"id": id}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (p *MongoDB) UnassignJobs(schedulerId string) (int64, error) {
	col := p.client.Database("spyglass").Collection("Jobs")

	res, err := col.UpdateMany(p.context, bson.M{"schedulerid": schedulerId},
		bson.D{{"$set", bson.M{"schedulerId": nil}}})
	if err != nil {
		p.Log.Error(err.Error())
		return 0, fmt.Errorf("Error updating jobs")
	}
	return res.MatchedCount, nil
}
