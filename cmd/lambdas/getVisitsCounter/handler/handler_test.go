package handler

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/jhidalgoesp/example-services/pkg/counter"
	"github.com/jhidalgoesp/example-services/pkg/web"
	"github.com/jhidalgoesp/example-services/tests"
	"net/http"
	"testing"
)

var counterMock = counter.Counter{
	ID:      "123",
	Counter: "20",
}

type CounterCoreMock struct {
	GetVisitsCounterCalls       *int
	IncrementVisitsCounterCalls *int
	error                       error
}

func (c CounterCoreMock) IncrementVisitsCounter() error {
	*c.IncrementVisitsCounterCalls++
	if c.error != nil {
		return c.error
	}

	return nil
}

func (c CounterCoreMock) GetVisitsCounter(counter *counter.Counter) error {
	*c.GetVisitsCounterCalls++
	if c.error != nil {
		return c.error
	}

	*counter = counterMock

	return nil
}

func NewCounterCoreMock(error error) CounterCoreMock {
	return CounterCoreMock{
		GetVisitsCounterCalls:       new(int),
		IncrementVisitsCounterCalls: new(int),
		error:                       error,
	}
}

func TestGetVisitsCounterHandler_Handle(t *testing.T) {
	t.Run("Handler should return a ApiGatewayResponse containing an Counter", func(t *testing.T) {
		handler := GetVisitsCounterHandler{
			Core: NewCounterCoreMock(nil),
		}

		response, err := handler.Handle(events.APIGatewayProxyRequest{})
		tests.AssertNilError(t, err)

		var counter counter.Counter

		err = json.Unmarshal([]byte(response.Body), &counter)

		tests.AssertNilError(t, err)
		tests.AssertStrings(t, counter.ID, counterMock.ID)
		tests.AssertStrings(t, counter.Counter, counterMock.Counter)
	})

	t.Run("Handler should return a status 200", func(t *testing.T) {
		handler := GetVisitsCounterHandler{
			Core: NewCounterCoreMock(nil),
		}
		response, err := handler.Handle(events.APIGatewayProxyRequest{})

		tests.AssertNilError(t, err)
		tests.AssertSameInt(t, response.StatusCode, http.StatusOK)
	})

	t.Run("Handler should call GetVisitsCounter once", func(t *testing.T) {
		core := NewCounterCoreMock(nil)
		handler := GetVisitsCounterHandler{
			Core: core,
		}
		_, err := handler.Handle(events.APIGatewayProxyRequest{})

		tests.AssertNilError(t, err)
		if *core.GetVisitsCounterCalls != 1 {
			t.Errorf("Expected GetVisitsCounterCalls to be called once got %d", *core.GetVisitsCounterCalls)
		}
	})

	t.Run("Handler should return a ErrorResponse if GetVisitsCounter fails", func(t *testing.T) {
		er := web.RequestError{
			Err:    fmt.Errorf("database error"),
			Status: http.StatusBadRequest,
		}
		handler := GetVisitsCounterHandler{
			Core: NewCounterCoreMock(er.Err),
		}

		_, err := handler.Handle(events.APIGatewayProxyRequest{})
		tests.AssertNotNilError(t, err)
		tests.AssertStrings(t, err.Error(), er.Err.Error())
	})
}
