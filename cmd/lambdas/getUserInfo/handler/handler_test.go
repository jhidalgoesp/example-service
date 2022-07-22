package handler

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/jhidalgoesp/example-services/pkg/counter"
	"github.com/jhidalgoesp/example-services/pkg/user"
	"github.com/jhidalgoesp/example-services/pkg/web"
	"github.com/jhidalgoesp/example-services/tests"
	"net/http"
	"testing"
)

var mockUser = user.User{
	ID:             "123",
	Name:           "Jose Hidalgo",
	WorkExperience: "5 years",
	TwitterHandle:  "jhidalgoesp",
}

type UserCoreMock struct {
	GetUserByIdCalls *int
	error            error
}

func (u UserCoreMock) GetUserById(userId string, userDb *user.User) error {
	*u.GetUserByIdCalls++

	if u.error != nil {
		return u.error
	}

	*userDb = mockUser

	return nil
}

func (u UserCoreMock) Update(user.User, string) (user.User, error) {
	*u.GetUserByIdCalls++

	if u.error != nil {
		return user.User{}, u.error
	}

	return user.User{}, nil
}

func NewUserCoreMock(error error) UserCoreMock {
	return UserCoreMock{
		GetUserByIdCalls: new(int),
		error:            error,
	}
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

func (c CounterCoreMock) GetVisitsCounter(*counter.Counter) error {
	*c.GetVisitsCounterCalls++
	if c.error != nil {
		return c.error
	}

	return nil
}

func NewCounterCoreMock(error error) CounterCoreMock {
	return CounterCoreMock{
		GetVisitsCounterCalls:       new(int),
		IncrementVisitsCounterCalls: new(int),
		error:                       error,
	}
}

func TestGetUserInfoHandler_Handle(t *testing.T) {
	t.Run("Handler should return a ApiGatewayResponse containing an user", func(t *testing.T) {
		handler := GetUserInfoHandler{
			Users:   NewUserCoreMock(nil),
			Counter: NewCounterCoreMock(nil),
		}

		response, err := handler.Handle(events.APIGatewayProxyRequest{})
		tests.AssertNilError(t, err)

		var user user.User

		err = json.Unmarshal([]byte(response.Body), &user)

		tests.AssertNilError(t, err)
		tests.AssertStrings(t, user.ID, mockUser.ID)
		tests.AssertStrings(t, user.Name, mockUser.Name)
		tests.AssertStrings(t, user.TwitterHandle, mockUser.TwitterHandle)
		tests.AssertStrings(t, user.WorkExperience, mockUser.WorkExperience)
	})

	t.Run("Handler should return a status 200", func(t *testing.T) {
		handler := GetUserInfoHandler{
			Users:   NewUserCoreMock(nil),
			Counter: NewCounterCoreMock(nil),
		}
		response, err := handler.Handle(events.APIGatewayProxyRequest{})

		tests.AssertNilError(t, err)
		tests.AssertSameInt(t, response.StatusCode, http.StatusOK)
	})

	t.Run("Handler should call GetUserById once", func(t *testing.T) {
		core := NewUserCoreMock(nil)
		handler := GetUserInfoHandler{
			Users:   core,
			Counter: NewCounterCoreMock(nil),
		}
		_, err := handler.Handle(events.APIGatewayProxyRequest{})

		tests.AssertNilError(t, err)
		if *core.GetUserByIdCalls != 1 {
			t.Errorf("Expected GetVisitsCounterCalls to be called once got %d", *core.GetUserByIdCalls)
		}
	})

	t.Run("Handler should call IncrementVisitsCounter once", func(t *testing.T) {
		core := NewCounterCoreMock(nil)
		handler := GetUserInfoHandler{
			Users:   NewUserCoreMock(nil),
			Counter: core,
		}
		_, err := handler.Handle(events.APIGatewayProxyRequest{})

		tests.AssertNilError(t, err)
		if *core.IncrementVisitsCounterCalls != 1 {
			t.Errorf("Expected GetVisitsCounterCalls to be called once got %d", *core.IncrementVisitsCounterCalls)
		}
	})

	t.Run("Handler should return a ErrorResponse if GetUserById fails", func(t *testing.T) {
		er := web.RequestError{
			Err:    fmt.Errorf("database error"),
			Status: http.StatusBadRequest,
		}
		handler := GetUserInfoHandler{
			Users:   NewUserCoreMock(fmt.Errorf("database error")),
			Counter: NewCounterCoreMock(nil),
		}

		_, err := handler.Handle(events.APIGatewayProxyRequest{})
		tests.AssertNotNilError(t, err)
		tests.AssertStrings(t, err.Error(), er.Err.Error())
	})

	t.Run("Handler should return a ErrorResponse if IncrementVisitsCounter fails", func(t *testing.T) {
		er := web.RequestError{
			Err:    fmt.Errorf("database error"),
			Status: http.StatusBadRequest,
		}
		handler := GetUserInfoHandler{
			Users:   NewUserCoreMock(nil),
			Counter: NewCounterCoreMock(fmt.Errorf("database error")),
		}

		_, err := handler.Handle(events.APIGatewayProxyRequest{})
		tests.AssertNotNilError(t, err)
		tests.AssertStrings(t, err.Error(), er.Err.Error())
	})
}
