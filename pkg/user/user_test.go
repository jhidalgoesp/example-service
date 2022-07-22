package user

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jhidalgoesp/example-services/tests"
	"testing"
)

type StoreMock struct {
	store
	returnGetItem *dynamodb.GetItemOutput
	error         error
}

func newStoreMock(error error, returnGetItem *dynamodb.GetItemOutput) StoreMock {
	return StoreMock{
		error:         error,
		returnGetItem: returnGetItem,
	}
}

func (s StoreMock) getItem(userID string) (*dynamodb.GetItemOutput, error) {
	if s.error != nil {
		return nil, s.error
	}

	if s.returnGetItem != nil {
		return s.returnGetItem, nil
	}

	return &dynamodb.GetItemOutput{
		Item: dynamoResultMock,
	}, nil
}

func (s StoreMock) updateItem(user User, userId string) (*dynamodb.UpdateItemOutput, error) {
	if s.error != nil {
		return nil, s.error
	}

	return &dynamodb.UpdateItemOutput{
		Attributes: dynamoResultMock,
	}, nil
}

func TestCore_GetUserById(t *testing.T) {
	t.Run("GetUserById should return a User", func(t *testing.T) {
		store := newStoreMock(nil, nil)
		core := Core{
			store: store,
		}

		user := User{}
		err := core.GetUserById("456", &user)

		tests.AssertNilError(t, err)
		tests.AssertStrings(t, user.ID, "123")
	})

	t.Run("GetUserById should return ErrInvalid if the userId param is empty", func(t *testing.T) {
		store := newStoreMock(nil, nil)
		core := Core{
			store: store,
		}

		user := User{}
		err := core.GetUserById("", &user)

		tests.AssertSameErrors(t, err, ErrInvalidId)
	})

	t.Run("GetUserById should return ErrDatabase if the query fails", func(t *testing.T) {
		store := newStoreMock(fmt.Errorf("query error"), nil)
		core := Core{
			store: store,
		}

		user := User{}
		err := core.GetUserById("123", &user)

		tests.AssertSameErrors(t, err, ErrDatabase)
	})

	t.Run("GetUserById should return ErrNotFound if the user is not in the store", func(t *testing.T) {
		store := newStoreMock(nil,
			&dynamodb.GetItemOutput{
				Item: nil,
			})

		core := Core{
			store: store,
		}

		user := User{}
		err := core.GetUserById("123", &user)

		tests.AssertSameErrors(t, err, ErrNotFound)
	})
}

func TestCore_Update(t *testing.T) {
	t.Run("Update should return a User", func(t *testing.T) {
		store := newStoreMock(nil, nil)
		core := Core{
			store: store,
		}

		result, err := core.Update(User{}, "123")

		tests.AssertNilError(t, err)
		tests.AssertStrings(t, result.ID, "123")
	})

	t.Run("Update should return ErrDatabase if the query fails", func(t *testing.T) {
		store := newStoreMock(fmt.Errorf("query error"), nil)
		core := Core{
			store: store,
		}

		_, err := core.Update(User{}, "123")

		tests.AssertSameErrors(t, err, ErrDatabase)
	})
}

func TestUser_NewCore(t *testing.T) {
	t.Run("NewCore should return a Core", func(t *testing.T) {
		core := NewCore(nil, nil)
		tests.AssertKindIsStruct(t, core)
	})
}
