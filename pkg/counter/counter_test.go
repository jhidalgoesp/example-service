package counter

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/jhidalgoesp/example-services/tests"
	"testing"
)

type StoreMock struct {
	store
	returnGetItem         *dynamodb.GetItemOutput
	error                 error
	incrementCounterCalls *int
	getCounterCalls       *int
}

func newStoreMock(error error, returnGetItem *dynamodb.GetItemOutput) StoreMock {
	return StoreMock{
		error:                 error,
		returnGetItem:         returnGetItem,
		incrementCounterCalls: new(int),
		getCounterCalls:       new(int),
	}
}

func (s StoreMock) getCounter() (*dynamodb.GetItemOutput, error) {
	*s.getCounterCalls++
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

func (s StoreMock) incrementCounter() error {
	*s.incrementCounterCalls++
	if s.error != nil {
		return s.error
	}

	return nil
}

func TestCounter_IncrementVisitsCounter(t *testing.T) {
	t.Run("IncrementVisitsCounter should return a nil error", func(t *testing.T) {
		store := newStoreMock(nil, nil)
		core := Core{store}
		err := core.IncrementVisitsCounter()
		tests.AssertNilError(t, err)
	})

	t.Run("incrementCounter should be called once.", func(t *testing.T) {
		store := newStoreMock(nil, nil)
		core := Core{store}

		err := core.IncrementVisitsCounter()

		tests.AssertNilError(t, err)
		if *store.incrementCounterCalls != 1 {
			t.Errorf("Expected getItem to be called 1 time got called %d times instead.",
				*store.incrementCounterCalls)
		}
	})

	t.Run("IncrementVisitsCounter should return an error if query fails.", func(t *testing.T) {
		store := newStoreMock(fmt.Errorf("database error"), nil)
		core := Core{store}

		err := core.IncrementVisitsCounter()

		tests.AssertNotNilError(t, err)
	})
}

func TestCounter_GetVisitsCounter(t *testing.T) {
	t.Run("GetVisitsCounter should return a nil error", func(t *testing.T) {
		store := newStoreMock(nil, nil)
		core := Core{store}

		counter := Counter{}
		err := core.GetVisitsCounter(&counter)

		tests.AssertNilError(t, err)
		tests.AssertStrings(t, counter.ID, "visitsCounter")
		tests.AssertStrings(t, counter.Counter, "123")
	})

	t.Run("getCounter should be called once.", func(t *testing.T) {
		store := newStoreMock(nil, nil)
		core := Core{store}

		counter := Counter{}
		err := core.GetVisitsCounter(&counter)

		tests.AssertNilError(t, err)
		if *store.getCounterCalls != 1 {
			t.Errorf("Expected getCounter to be called 1 time got called %d times instead.",
				*store.getCounterCalls)
		}
	})

	t.Run("GetCounter should return an error if query fails.", func(t *testing.T) {
		store := newStoreMock(fmt.Errorf("database error"), nil)
		core := Core{store}

		counter := Counter{}
		err := core.GetVisitsCounter(&counter)

		tests.AssertNotNilError(t, err)
	})

	t.Run("GetCounter should return ErrNotFound if the key is not in the store.", func(t *testing.T) {
		store := newStoreMock(nil, &dynamodb.GetItemOutput{
			Item: nil,
		})
		core := Core{store}

		counter := Counter{}
		err := core.GetVisitsCounter(&counter)

		tests.AssertSameErrors(t, err, ErrNotFound)
	})
}

func TestCounter_NewCore(t *testing.T) {
	t.Run("NewCore should return a Core", func(t *testing.T) {
		core := NewCore(nil, nil)
		tests.AssertKindIsStruct(t, core)
	})
}
