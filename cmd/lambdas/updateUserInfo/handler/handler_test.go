package handler

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
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
	UpdateCalls      *int
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
	*u.UpdateCalls++

	if u.error != nil {
		return user.User{}, u.error
	}

	return mockUser, nil
}

func NewUserCoreMock(error error) UserCoreMock {
	return UserCoreMock{
		GetUserByIdCalls: new(int),
		UpdateCalls:      new(int),
		error:            error,
	}
}

func TestUpdateUserInfoHandler_Handle(t *testing.T) {
	t.Run("Handler should return a ApiGatewayResponse containing a User", func(t *testing.T) {
		handler := UpdateUserInfoHandler{
			Users: NewUserCoreMock(nil),
		}
		userJson, err := json.Marshal(mockUser)
		tests.AssertNilError(t, err)

		response, err := handler.Handle(events.APIGatewayProxyRequest{Body: string(userJson)})
		tests.AssertNilError(t, err)
		var userDb user.User

		err = json.Unmarshal([]byte(response.Body), &userDb)
		tests.AssertNilError(t, err)
		tests.AssertStrings(t, userDb.ID, mockUser.ID)
		tests.AssertStrings(t, userDb.Name, mockUser.Name)
		tests.AssertStrings(t, userDb.TwitterHandle, mockUser.TwitterHandle)
		tests.AssertStrings(t, userDb.WorkExperience, mockUser.WorkExperience)
	})

	t.Run("Handler should return a status 200", func(t *testing.T) {
		handler := UpdateUserInfoHandler{
			Users: NewUserCoreMock(nil),
		}
		userJson, err := json.Marshal(mockUser)
		tests.AssertNilError(t, err)

		response, err := handler.Handle(events.APIGatewayProxyRequest{Body: string(userJson)})

		tests.AssertNilError(t, err)
		tests.AssertSameInt(t, response.StatusCode, http.StatusOK)
	})

	t.Run("Handler should call user.Update once", func(t *testing.T) {
		core := NewUserCoreMock(nil)
		handler := UpdateUserInfoHandler{
			Users: core,
		}
		userJson, err := json.Marshal(mockUser)
		tests.AssertNilError(t, err)

		_, err = handler.Handle(events.APIGatewayProxyRequest{Body: string(userJson)})

		tests.AssertNilError(t, err)
		if *core.UpdateCalls != 1 {
			t.Errorf("Expected GetVisitsCounterCalls to be called once got %d", *core.UpdateCalls)
		}
	})

	t.Run("Handler should return a ErrorResponse if User.Update fails", func(t *testing.T) {
		er := web.RequestError{
			Err:    fmt.Errorf("database error"),
			Status: http.StatusBadRequest,
		}
		handler := UpdateUserInfoHandler{Users: NewUserCoreMock(fmt.Errorf("database error"))}
		userJson, err := json.Marshal(mockUser)
		tests.AssertNilError(t, err)

		_, err = handler.Handle(events.APIGatewayProxyRequest{Body: string(userJson)})
		tests.AssertNotNilError(t, err)
		tests.AssertStrings(t, err.Error(), er.Err.Error())
	})

	t.Run("Handler should return a ErrorResponse if request does not contains an user fails", func(t *testing.T) {
		handler := UpdateUserInfoHandler{Users: NewUserCoreMock(fmt.Errorf("database error"))}
		_, err := handler.Handle(events.APIGatewayProxyRequest{})
		tests.AssertNotNilError(t, err)
	})
}
