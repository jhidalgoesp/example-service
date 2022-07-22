package handler

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/jhidalgoesp/example-services/pkg/twitter"
	"github.com/jhidalgoesp/example-services/pkg/web"
	"github.com/jhidalgoesp/example-services/tests"
	"net/http"
	"testing"
)

type TwitterClientMock struct {
	GetLastTweetsCalls *int
	error              error
}

func (t TwitterClientMock) GetLastTweets(username, count string, tweets *[]twitter.Tweet) error {
	*t.GetLastTweetsCalls++
	if t.error != nil {
		return t.error
	}

	*tweets = []twitter.Tweet{
		{
			Id:        123,
			CreatedAt: "now",
			Text:      "This a tweet mock",
			User: twitter.User{
				Name:         "Jose Hidalgo",
				ScreenName:   "jhidalgoesp",
				ProfileImage: "http://test.com",
			},
		},
		{
			Id:        234,
			CreatedAt: "now",
			Text:      "This a tweet mock",
			User: twitter.User{
				Name:         "Mariana Rodas",
				ScreenName:   "arodas",
				ProfileImage: "http://test.com",
			},
		},
	}

	return nil
}

func NewTwitterClientMock(error error) TwitterClientMock {
	return TwitterClientMock{
		GetLastTweetsCalls: new(int),
		error:              error,
	}
}

func TestLastTweetsHandler_Handle(t *testing.T) {
	t.Run("Handler should return a ApiGatewayResponse containing two tweets", func(t *testing.T) {
		handler := LastTweetsHandler{Client: NewTwitterClientMock(nil)}

		response, err := handler.Handle(events.APIGatewayProxyRequest{})

		tests.AssertNilError(t, err)

		var tweets []twitter.Tweet

		err = json.Unmarshal([]byte(response.Body), &tweets)

		tests.AssertNilError(t, err)
		tests.AssertSliceLength(t, len(tweets), 2)
	})

	t.Run("Handler should return a status 200", func(t *testing.T) {
		handler := LastTweetsHandler{Client: NewTwitterClientMock(nil)}

		response, err := handler.Handle(events.APIGatewayProxyRequest{})

		tests.AssertNilError(t, err)
		tests.AssertSameInt(t, response.StatusCode, http.StatusOK)
	})

	t.Run("Handler should return a ErrorResponse if GetLastTweets fails", func(t *testing.T) {
		er := web.RequestError{
			Err:    fmt.Errorf("database error"),
			Status: http.StatusBadRequest,
		}
		handler := LastTweetsHandler{Client: NewTwitterClientMock(fmt.Errorf("database error"))}

		_, err := handler.Handle(events.APIGatewayProxyRequest{})
		tests.AssertNotNilError(t, err)
		tests.AssertStrings(t, err.Error(), er.Err.Error())
	})
}
