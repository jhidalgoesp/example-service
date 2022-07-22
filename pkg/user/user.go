// Package user provides the client and interface to interact with users.
package user

import (
	"errors"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/jhidalgoesp/example-services/pkg/validate"
	"go.uber.org/zap"
)

var (
	ErrInvalidId   = errors.New("id is not valid")
	ErrNotFound    = errors.New("user not found")
	ErrDatabase    = errors.New("user not found")
	ErrMarshalling = errors.New("user not found")
)

// CoreAPI describes users business operations.
type CoreAPI interface {
	GetUserById(userId string, userDb *User) error
	Update(User, string) (User, error)
}

// User represents an individual user.
type User struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	WorkExperience string `json:"workExperience"`
	TwitterHandle  string `json:"twitterHandle"`
}

// Core manages the set of APIs for user access.
type Core struct {
	store store
}

// NewCore constructs a core for user api access.
func NewCore(log *zap.SugaredLogger, dynamo *dynamodb.DynamoDB) Core {
	return Core{
		store: newStore(log, dynamo),
	}
}

// GetUserById gets the specified user from the dynamoStore.
func (c Core) GetUserById(userId string, userDb *User) error {
	if validate.IsEmpty(userId) {
		return ErrInvalidId
	}

	result, err := c.store.getItem(userId)

	if err != nil {
		return ErrDatabase
	}

	if result.Item == nil {
		return ErrNotFound
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &userDb)
	if err != nil {
		return ErrMarshalling
	}

	return nil
}

// Update updates a specific record from the dynamoStore.
func (c Core) Update(user User, userId string) (User, error) {
	dbUser := User{}

	result, err := c.store.updateItem(user, userId)

	if err != nil {
		return dbUser, ErrDatabase
	}

	err = dynamodbattribute.UnmarshalMap(result.Attributes, &dbUser)
	if err != nil {
		return dbUser, ErrMarshalling
	}

	return dbUser, nil
}
