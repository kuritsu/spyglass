package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/google/uuid"
	"github.com/kuritsu/spyglass/api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	Log         *logrus.Logger
	initialized bool
	client      *mongo.Client
	context     context.Context
	cancelFunc  context.CancelFunc
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

	if p.initialized {
		return
	}

	// Force a connection to verify our connection string
	err = client.Ping(ctx, nil)
	if err != nil {
		p.Log.Fatalf("Failed to ping cluster: %v", err)
	}

	p.Log.Info("Connected to MongoDB!")

	p.createIndexes()
	p.initialized = true
}

func (p *MongoDB) createIndexes() {
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
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	_, err := roleCollection.InsertOne(p.context, adminsRole)
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
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
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

func (p *MongoDB) GetAllUsers(pageSize, pageIndex int64) ([]*types.User, error) {
	col := p.client.Database("spyglass").Collection("Users")
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
	users := make([]*types.User, 0)
	err = cursor.All(p.context, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
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

// InsertMonitor into the db
func (p *MongoDB) InsertMonitor(monitor *types.Monitor) (*types.Monitor, error) {
	monitor.CreatedAt = time.Now()
	monitor.UpdatedAt = time.Now()
	_, err := p.client.Database("spyglass").Collection("Monitors").InsertOne(
		p.context, monitor)
	if err != nil {
		p.Log.Errorf("Could not create Monitor: %v", err)
		return nil, err
	}
	return monitor, nil
}

// InsertTarget into the db.
func (p *MongoDB) InsertTarget(target *types.Target) (*types.Target, error) {
	target.ID = strings.ToLower(target.ID)
	col := p.client.Database("spyglass").Collection("Targets")
	allTargets := updateChildrenRefs(target)
	p.Log.Debug("Found ", len(allTargets), " targets.")
	allObjects := make([]interface{}, len(allTargets))
	for i := range allTargets {
		allTargets[i].CreatedAt = time.Now()
		allTargets[i].UpdatedAt = time.Now()
		allObjects[i] = allTargets[i]
	}
	p.Log.Debug("Calling InsertMany...")
	if _, err := col.InsertMany(p.context, allObjects, options.InsertMany()); err != nil {
		p.Log.Errorf("Could not create Target: %v", err)
		return nil, err
	}
	return target, nil
}

// UpdateMonitor with the modifyable fields.
func (p *MongoDB) UpdateMonitor(oldMonitor *types.Monitor, newMonitor *types.Monitor) (*types.Monitor, error) {
	newMonitor.CreatedAt = oldMonitor.CreatedAt
	newMonitor.Owners = oldMonitor.Owners
	newMonitor.UpdatedAt = time.Now()
	_, err := p.client.Database("spyglass").Collection("Monitors").UpdateOne(
		p.context, bson.M{"id": oldMonitor.ID},
		bson.M{"$set": newMonitor})
	if err != nil {
		p.Log.Errorf("Could not create Monitor: %v", err)
		return nil, err
	}
	return newMonitor, nil
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
	newTarget.Permissions.UpdatedAt = time.Now()
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

func (p *MongoDB) Login(email string, password string) (*types.User, error) {
	col := p.client.Database("spyglass").Collection("Users")
	expr := bson.M{"email": email}
	res := col.FindOne(p.context, expr)
	if res.Err() != nil {
		p.Log.Error(res.Err())
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("InvalidCredentials")
		} else {
			return nil, res.Err()
		}
	}
	var user types.User
	res.Decode(&user)
	err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(password))
	if err != nil {
		return nil, fmt.Errorf("InvalidCredentials")
	}
	return &user, nil
}

func (p *MongoDB) Register(email string, password string) (*types.User, error) {
	epwd, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	user := types.User{
		Email:     email,
		FullName:  email,
		Roles:     []string{email},
		PassHash:  string(epwd),
		FirstHash: string(epwd),
		Permissions: types.Permissions{
			Owners:    []string{email},
			Writers:   []string{email},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}
	_, err := p.client.Database("spyglass").Collection("Users").InsertOne(p.context, user)
	if err != nil {
		p.Log.Error(err)
		return nil, fmt.Errorf("ErrorCreatingUser")
	}
	return &user, nil
}

func (p *MongoDB) CreateUserToken(user *types.User) (string, error) {
	tokenUuid := uuid.NewString()
	token := types.UserToken{
		Email:      user.Email,
		Expiration: time.Now().UTC().Add(time.Hour * 24),
		Token:      tokenUuid,
	}
	_, err := p.client.Database("spyglass").Collection("Tokens").InsertOne(
		p.context, token)
	if err != nil {
		return "", err
	}
	return tokenUuid, nil
}

func (p *MongoDB) ValidateToken(email string, token string) error {
	col := p.client.Database("spyglass").Collection("Tokens")
	expr := bson.M{"email": email, "token": token}
	res := col.FindOne(p.context, expr)
	if res.Err() != nil {
		p.Log.Error(res.Err())
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return fmt.Errorf("InvalidCredentials")
		} else {
			return res.Err()
		}
	}
	var t types.UserToken
	err := res.Decode(&t)
	if err == nil && t.Expiration.After(time.Now().UTC()) {
		return nil
	}
	return fmt.Errorf("InvalidCredentials")
}

func (p *MongoDB) GetUser(email string) (*types.User, error) {
	expr := bson.M{"email": email}
	res := p.client.Database("spyglass").Collection("Users").FindOne(p.context, expr)
	if res.Err() != nil {
		p.Log.Error(res.Err())
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return nil, fmt.Errorf("InvalidUser")
		} else {
			return nil, res.Err()
		}
	}
	var user types.User
	err := res.Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (p *MongoDB) UpdateUser(user *types.User, oldPassword, newPassword string) error {
	if oldPassword != "" {
		err := bcrypt.CompareHashAndPassword([]byte(user.PassHash), []byte(oldPassword))
		if err != nil {
			return fmt.Errorf("Invalid credentials")
		}
		epwd, err := bcrypt.GenerateFromPassword([]byte(newPassword), 14)
		if err != nil {
			return fmt.Errorf("Encryption Error")
		}
		user.PassHash = string(epwd)
	}
	user.Permissions.UpdatedAt = time.Now()
	_, err := p.client.Database("spyglass").Collection("Users").UpdateOne(
		p.context, bson.M{"email": user.Email},
		bson.M{"$set": user})
	if err != nil {
		p.Log.Errorf("Could not update user: %v", err)
		return err
	}
	return nil

}

func (p *MongoDB) InsertRole(role *types.Role, user *types.User) error {
	role.CreatedAt = time.Now()
	role.UpdatedAt = time.Now()

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
