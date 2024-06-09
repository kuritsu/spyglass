package mongodb

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kuritsu/spyglass/api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

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
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		},
	}
	_, err := p.client.Database("spyglass").Collection("Users").InsertOne(p.context, user)
	if err != nil {
		p.Log.Error(err)
		return nil, fmt.Errorf("ErrorCreatingUser")
	}
	return &user, nil
}

func (p *MongoDB) CreateUserToken(user *types.User, expiration time.Time) (string, error) {
	tokenUuid := uuid.NewString()
	token := types.UserToken{
		Email:      user.Email,
		Expiration: expiration,
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
	user.Permissions.UpdatedAt = time.Now().UTC()
	_, err := p.client.Database("spyglass").Collection("Users").UpdateOne(
		p.context, bson.M{"email": user.Email},
		bson.M{"$set": user})
	if err != nil {
		p.Log.Errorf("Could not update user: %v", err)
		return err
	}
	return nil
}
