package counter

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jhidalgoesp/example-services/tests"
	"testing"
)

var dynamoResultMock = map[string]*dynamodb.AttributeValue{
	"Id":      {S: aws.String("visitsCounter")},
	"Counter": {N: aws.String("123")},
}

func TestStore_IncrementCounter(t *testing.T) {
	t.Run("incrementCounter should execute dynamodb UpdateItem once", func(t *testing.T) {
		dynamo := tests.NewDynamoMock(nil, nil, &dynamodb.UpdateItemOutput{
			Attributes: dynamoResultMock,
		})

		store := dynamoStore{
			dynamo: dynamo,
			log:    nil,
		}

		err := store.incrementCounter()

		tests.AssertNilError(t, err)

		if *dynamo.UpdateItemCount != 1 {
			t.Errorf("Expected UpdateItemCount to be called once, got called %d times.", *dynamo.UpdateItemCount)
		}
	})

	t.Run("incrementCounter should return an error if the query fails", func(t *testing.T) {
		dynamo := tests.NewDynamoMock(fmt.Errorf("database error"), nil, nil)

		store := dynamoStore{
			dynamo: dynamo,
			log:    nil,
		}

		err := store.incrementCounter()

		tests.AssertNotNilError(t, err)
	})
}

func TestStore_GetCounter(t *testing.T) {
	t.Run("getCounter should return a GetItemOutput pointer", func(t *testing.T) {
		dynamo := tests.NewDynamoMock(nil, &dynamodb.GetItemOutput{
			Item: dynamoResultMock,
		}, nil)

		store := dynamoStore{
			dynamo: dynamo,
			log:    nil,
		}

		result, err := store.getCounter()

		tests.AssertNilError(t, err)
		tests.AssertStrings(t, *result.Item["Id"].S, "visitsCounter")
		tests.AssertStrings(t, *result.Item["Counter"].N, "123")
	})

	t.Run("getCounter should execute dynamodb GetItem once", func(t *testing.T) {
		dynamo := tests.NewDynamoMock(nil, &dynamodb.GetItemOutput{
			Item: dynamoResultMock,
		}, nil)

		store := dynamoStore{
			dynamo: dynamo,
			log:    nil,
		}

		_, err := store.getCounter()

		tests.AssertNilError(t, err)

		if *dynamo.GetItemCount != 1 {
			t.Errorf("Expected UpdateItemCount to be called once, got called %d times.", *dynamo.UpdateItemCount)
		}
	})

	t.Run("getCounter should return an error if the query fails", func(t *testing.T) {
		dynamo := tests.NewDynamoMock(fmt.Errorf("database error"), nil, nil)

		store := dynamoStore{
			dynamo: dynamo,
			log:    nil,
		}

		_, err := store.getCounter()
		tests.AssertNotNilError(t, err)
	})
}

func TestStore_NewStore(t *testing.T) {
	t.Run("NewStore should return a Store", func(t *testing.T) {
		store := NewStore(nil, nil)
		tests.AssertKindIsStruct(t, store)
	})
}
