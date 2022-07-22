// Package handler provides the http logic of the service.
package handler

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/jhidalgoesp/example-services/pkg/counter"
	"github.com/jhidalgoesp/example-services/pkg/user"
	"github.com/jhidalgoesp/example-services/pkg/web"
	"go.uber.org/zap"
	"net/http"
)

// GetUserInfoHandler contains all the requirements of the handler
type GetUserInfoHandler struct {
	Log     *zap.SugaredLogger
	Users   user.CoreAPI
	Counter counter.CoreAPI
}

func (h GetUserInfoHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userId := request.QueryStringParameters["id"]

	user := user.User{}
	err := h.Users.GetUserById(userId, &user)

	if err != nil {
		return events.APIGatewayProxyResponse{}, web.NewRequestError(err, http.StatusBadRequest)
	}

	err = h.Counter.IncrementVisitsCounter()
	if err != nil {
		return events.APIGatewayProxyResponse{}, web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Respond(user, http.StatusOK)
}
