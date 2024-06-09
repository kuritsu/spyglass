package mongodb

import (
	"errors"
	"fmt"
	"time"

	"github.com/kuritsu/spyglass/api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (p *MongoDB) GetAllRoles(pageSize, pageIndex int64) ([]*types.Role, error) {
	col := p.client.Database("spyglass").Collection("Roles")
	skip := pageIndex * pageSize
	opts := options.FindOptions{
		Skip:  &skip,
		Sort:  bson.M{"id": 1},
		Limit: &pageSize,
	}
	cursor, err := col.Find(p.context, bson.D{}, &opts)
	if err != nil {
		return nil, err
	}
	roles := make([]*types.Role, 0)
	err = cursor.All(p.context, &roles)
	if err != nil {
		return nil, err
	}
	return roles, nil
}

func (p *MongoDB) InsertRole(role *types.Role, user *types.User) error {
	role.CreatedAt = time.Now().UTC()
	role.UpdatedAt = time.Now().UTC()

	_, err := p.client.Database("spyglass").Collection("Roles").InsertOne(p.context, role)
	if err != nil {
		p.Log.Error(err)
		return fmt.Errorf("Error creating role")
	}
	return nil
}

func (p *MongoDB) GetRole(name string) (*types.Role, error) {
	expr := bson.M{"name": name}
	res := p.client.Database("spyglass").Collection("Roles").FindOne(p.context, expr)
	if res.Err() != nil {
		p.Log.Error(res.Err())
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("InvalidRole")
		} else {
			return nil, res.Err()
		}
	}
	var role types.Role
	err := res.Decode(&role)
	if err != nil {
		return nil, err
	}
	return &role, nil
}
