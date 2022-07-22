package user

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	"go.uber.org/zap"
)

// store describes datastore access behavior.
type store interface {
	getItem(userID string) (*dynamodb.GetItemOutput, error)
	updateItem(user User, userId string) (*dynamodb.UpdateItemOutput, error)
}

// dynamoStore manages the set of APIs for user dynamodb access.
type dynamoStore struct {
	dynamo dynamodbiface.DynamoDBAPI
	log    *zap.SugaredLogger
}

// newStore constructs a data for api access.
func newStore(log *zap.SugaredLogger, dynamo *dynamodb.DynamoDB) dynamoStore {
	return dynamoStore{
		log:    log,
		dynamo: dynamo,
	}
}

// getItem retrieves a record from the Users table.
func (s dynamoStore) getItem(userID string) (*dynamodb.GetItemOutput, error) {
	result, err := s.dynamo.GetItem(&dynamodb.GetItemInput{
		Key:       map[string]*dynamodb.AttributeValue{"Id": {S: aws.String(userID)}},
		TableName: aws.String("test-services-users"),
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// updateItem updates a record form the Users table.
func (s dynamoStore) updateItem(user User, userId string) (*dynamodb.UpdateItemOutput, error) {
	result, err := s.dynamo.UpdateItem(&dynamodb.UpdateItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"Id": {S: aws.String(userId)},
		},
		TableName: aws.String("test-services-users"),
		UpdateExpression: aws.String(
			"SET #name = :name, " +
				"#workExperience = :workExperience, " +
				"#twitterHandle = :twitterHandle",
		),
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name":           {S: aws.String(user.Name)},
			":workExperience": {S: aws.String(user.WorkExperience)},
			":twitterHandle":  {S: aws.String(user.TwitterHandle)},
		},
		ExpressionAttributeNames: map[string]*string{
			"#name":           aws.String("Name"),
			"#workExperience": aws.String("WorkExperience"),
			"#twitterHandle":  aws.String("TwitterHandle"),
		},
		ReturnValues: aws.String(dynamodb.ReturnValueAllNew),
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}
