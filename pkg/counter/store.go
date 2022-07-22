package counter

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"go.uber.org/zap"
)

// store describes datastore access behavior.
type store interface {
	incrementCounter() error
	getCounter() (*dynamodb.GetItemOutput, error)
}

// dynamoStore manages the set of APIs for counter dynamodb access.
type dynamoStore struct {
	dynamo dynamodbiface.DynamoDBAPI
	log    *zap.SugaredLogger
}

// NewStore constructs a data for api access.
func NewStore(log *zap.SugaredLogger, dynamo *dynamodb.DynamoDB) dynamoStore {
	return dynamoStore{
		log:    log,
		dynamo: dynamo,
	}
}

// incrementCounter increments by one the profile visits atomic counter.
func (s dynamoStore) incrementCounter() error {
	_, err := s.dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		Key:                       map[string]*dynamodb.AttributeValue{"Id": {S: aws.String("userProfileVisits")}},
		TableName:                 aws.String("test-services-atomic-counter"),
		UpdateExpression:          aws.String("ADD #counter :val"),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{":val": {N: aws.String("1")}},
		ExpressionAttributeNames:  map[string]*string{"#counter": aws.String("Counter")},
	})

	if err != nil {
		return err
	}

	return nil
}

// getCounter retrieves the atomic counter for profile visits.
func (s dynamoStore) getCounter() (*dynamodb.GetItemOutput, error) {
	result, err := s.dynamo.GetItem(&dynamodb.GetItemInput{
		Key:       map[string]*dynamodb.AttributeValue{"Id": {S: aws.String("userProfileVisits")}},
		TableName: aws.String("test-services-atomic-counter"),
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
