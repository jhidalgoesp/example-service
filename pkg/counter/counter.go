// Package counter provides the client and interface to interact with atomic counters.
package counter

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"go.uber.org/zap"
)

var (
	ErrNotFound    = errors.New("counter not found")
	ErrDatabase    = errors.New("user not found")
	ErrMarshalling = errors.New("user not found")
)

// CoreAPI describes atomic counter operations.
type CoreAPI interface {
	IncrementVisitsCounter() error
	GetVisitsCounter(counter *Counter) error
}

// Counter represents a single record on the atomic counter table.
type Counter struct {
	ID      string `json:"id"`
	Counter string `json:"counter"`
}

// Core manages the set of APIs for atomic counter access.
type Core struct {
	store store
}

// NewCore constructs a core for atomic counter api access.
func NewCore(log *zap.SugaredLogger, dynamo *dynamodb.DynamoDB) Core {
	return Core{
		store: NewStore(log, dynamo),
	}
}

// IncrementVisitsCounter increments the visits counter by one.
func (c Core) IncrementVisitsCounter() error {
	err := c.store.incrementCounter()

	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	return nil
}

// GetVisitsCounter gets the visits counter from the store.
func (c Core) GetVisitsCounter(counter *Counter) error {
	result, err := c.store.getCounter()

	if err != nil {
		return ErrDatabase
	}

	if result.Item == nil {
		return ErrNotFound
	}

	err = dynamodbattribute.UnmarshalMap(result.Item, &counter)
	if err != nil {
		return ErrMarshalling
	}

	return nil
}
