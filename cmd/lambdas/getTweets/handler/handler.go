// Package handler provides the http logic of the service.
package handler

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/jhidalgoesp/example-services/pkg/twitter"
	"github.com/jhidalgoesp/example-services/pkg/web"
	"go.uber.org/zap"
	"net/http"
)

// LastTweetsHandler contains all the requirements of the handler
type LastTweetsHandler struct {
	Log    *zap.SugaredLogger
	Client twitter.Api
}

func (h LastTweetsHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	username := request.QueryStringParameters["username"]
	count := request.QueryStringParameters["count"]

	var tweets []twitter.Tweet
	err := h.Client.GetLastTweets(username, count, &tweets)

	if err != nil {
		return events.APIGatewayProxyResponse{}, web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Respond(tweets, http.StatusOK)
}
