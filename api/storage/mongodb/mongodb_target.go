package mongodb

import (
	"fmt"
	"strings"
	"time"

	"github.com/kuritsu/spyglass/api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GetAllTargets returns all targets which contains a string, paginated.
func (p *MongoDB) GetAllTargets(pageSize int64, pageIndex int64, contains string) ([]*types.Target, error) {
	col := p.client.Database("spyglass").Collection("Targets")
	containsRegex := `^[\w\d\-_\.]+$`
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
	return targets, nil
}

// GetTargetByID returns a target by its ID.
func (p *MongoDB) GetTargetByID(id string, includeChildren bool) (*types.Target, error) {
	col := p.client.Database("spyglass").Collection("Targets")
	id = strings.ToLower(id)
	expr := bson.M{"id": id}
	if includeChildren {
		escapedID := types.GetIDForRegex(id)
		regx := fmt.Sprintf(`^%s(/[\w\d_\-\.]+){0,1}$`, escapedID)
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
	var parent *types.Target
	var children []types.TargetRef
	for _, t := range targets {
		if t.ID == id {
			parent = t
			continue
		}
		children = append(children, t)
	}
	parent.Children = children
	if includeChildren {
		parent.ChildrenRef = nil
	}
	return parent, err
}

// InsertTarget into the db.
func (p *MongoDB) InsertTarget(target *types.Target) (*types.Target, error) {
	target.ID = strings.ToLower(target.ID)
	col := p.client.Database("spyglass").Collection("Targets")
	allTargets := updateChildrenRefs(target)
	p.Log.Debug("Found ", len(allTargets), " targets.")
	allObjects := make([]interface{}, len(allTargets))
	for i := range allTargets {
		allTargets[i].CreatedAt = time.Now().UTC()
		allTargets[i].UpdatedAt = time.Now().UTC()
		allObjects[i] = allTargets[i]
	}
	p.Log.Debug("Calling InsertMany...")
	if _, err := col.InsertMany(p.context, allObjects, options.InsertMany()); err != nil {
		p.Log.Errorf("Could not create Target: %v", err)
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
	target.Status = targetPatch.Status
	target.StatusDescription = targetPatch.StatusDescription
	return target, nil
}

// UpdateTarget with also a status update force flag.
func (p *MongoDB) UpdateTarget(oldTarget *types.Target, newTarget *types.Target,
	forceStatusUpdate bool) (*types.Target, error) {
	col := p.client.Database("spyglass").Collection("Targets")
	newTarget.Permissions.CreatedAt = oldTarget.CreatedAt
	newTarget.Permissions.UpdatedAt = time.Now().UTC()
	newTarget.Permissions.Owners = oldTarget.Owners
	updateProps := bson.M{
		"critical":    newTarget.Critical,
		"description": newTarget.Description,
		"monitor":     newTarget.Monitor,
		"permissions": newTarget.Permissions,
		"url":         newTarget.URL,
	}
	if forceStatusUpdate {
		lockedDoc := col.FindOneAndUpdate(p.context, bson.M{"id": oldTarget.ID},
			bson.M{"$set": bson.M{"flock": bson.M{"pseudoRandom": primitive.NewObjectID()}}})
		if lockedDoc.Err() != nil {
			return nil, lockedDoc.Err()
		}
		lockedDoc.Decode(oldTarget)
		updateProps["status"] = newTarget.Status
		updateProps["statusDescription"] = newTarget.StatusDescription
	}
	_, err := col.UpdateOne(p.context, bson.M{"id": oldTarget.ID},
		bson.M{"$set": updateProps})
	if err != nil {
		return nil, err
	}
	return newTarget, nil
}

func (p *MongoDB) DeleteTarget(id string) (int, error) {
	col := p.client.Database("spyglass").Collection("Targets")
	containsRegex := types.GetIDForRegex(id)
	containsRegex = fmt.Sprintf("%s(/.*)*", containsRegex)
	result, err := col.DeleteMany(p.context, bson.M{"id": bson.M{"$regex": containsRegex}}, nil)
	if err != nil {
		return 0, err
	}
	return int(result.DeletedCount), nil
}

func updateChildrenRefs(t *types.Target) []types.TargetRef {
	result := []types.TargetRef{t}
	if len(t.Children) == 0 {
		return result
	}
	var childrenRef []string = make([]string, len(t.Children))
	totalProgress := 0
	for i, c := range t.Children {
		c.ID = fmt.Sprintf("%s/%s", t.ID, strings.ToLower(c.ID))
		totalProgress += c.Status
		childrenRef[i] = types.GetShortID(c.ID)
		result = append(result, updateChildrenRefs(c)...)
	}
	t.Status = int(float64(100*totalProgress) / float64(100*len(t.Children)))
	t.ChildrenRef = childrenRef
	t.Children = nil
	return result
}
