package user

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jhidalgoesp/example-services/tests"
	"reflect"
	"testing"
)

var dynamoResultMock = map[string]*dynamodb.AttributeValue{
	"Id":             {S: aws.String("123")},
	"Name":           {S: aws.String("Jose Hidalgo")},
	"WorkExperience": {S: aws.String("Senior developer")},
	"TwitterHandle":  {S: aws.String("jhidalgoesp")},
}

func TestStore_GetItem(t *testing.T) {
	t.Run("GetItem should return a GetItemOutput", func(t *testing.T) {
		dynamo := tests.NewDynamoMock(nil, &dynamodb.GetItemOutput{
			Item: dynamoResultMock,
		}, nil)

		store := dynamoStore{
			dynamo: dynamo,
			log:    nil,
		}

		result, err := store.getItem("456")

		tests.AssertNilError(t, err)
		tests.AssertStrings(t, *result.Item["Id"].S, "123")
	})

	t.Run("GetItem should be called once.", func(t *testing.T) {
		dynamo := tests.NewDynamoMock(nil, &dynamodb.GetItemOutput{
			Item: dynamoResultMock,
		}, nil)
		store := dynamoStore{
			dynamo: dynamo,
			log:    nil,
		}

		_, err := store.getItem("123")

		tests.AssertNilError(t, err)
		if *dynamo.GetItemCount != 1 {
			t.Errorf("Expected getItem to be called 1 time got called %d times instead.",
				*dynamo.GetItemCount)
		}
	})

	t.Run("GetItem should return an error if query fails.", func(t *testing.T) {
		dynamo := tests.NewDynamoMock(fmt.Errorf("database error"), nil, nil)
		store := dynamoStore{
			dynamo: dynamo,
			log:    nil,
		}

		_, err := store.getItem("123")

		tests.AssertNotNilError(t, err)
	})
}

func TestStore_UpdateItem(t *testing.T) {
	t.Run("UpdateItem should return a UpdateItemOutput.", func(t *testing.T) {
		dynamo := tests.NewDynamoMock(nil, nil, &dynamodb.UpdateItemOutput{
			Attributes: dynamoResultMock,
		})
		store := dynamoStore{
			dynamo: dynamo,
			log:    nil,
		}

		result, err := store.updateItem(User{}, "123")

		tests.AssertNilError(t, err)
		tests.AssertStrings(t, *result.Attributes["Id"].S, "123")
	})

	t.Run("UpdateItem should be called one time.", func(t *testing.T) {
		dynamo := tests.NewDynamoMock(nil, nil, &dynamodb.UpdateItemOutput{
			Attributes: dynamoResultMock,
		})
		store := dynamoStore{
			dynamo: dynamo,
			log:    nil,
		}

		_, err := store.updateItem(User{}, "123")

		tests.AssertNilError(t, err)
		if *dynamo.UpdateItemCount != 1 {
			t.Errorf("Expected getItem to be called 1 time got called %d times instead.", *dynamo.GetItemCount)
		}
	})

	t.Run("UpdateItem should return an error if query fails.", func(t *testing.T) {
		dynamo := tests.NewDynamoMock(fmt.Errorf("database error"), nil, nil)
		store := dynamoStore{
			dynamo: dynamo,
			log:    nil,
		}

		_, err := store.updateItem(User{}, "123")

		tests.AssertNotNilError(t, err)
	})
}

func TestCore_NewStore(t *testing.T) {
	t.Run("NewStore should return a Store", func(t *testing.T) {
		store := newStore(nil, nil)
		storeType := reflect.TypeOf(store).Kind()

		if storeType != reflect.Struct {
			t.Errorf("Expected a user store struct, got %v", storeType)
		}
	})
}
