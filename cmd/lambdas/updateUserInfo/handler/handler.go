// Package handler provides the http logic of the service.
package handler

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/jhidalgoesp/example-services/pkg/user"
	"github.com/jhidalgoesp/example-services/pkg/web"
	"go.uber.org/zap"
	"net/http"
)

// UpdateUserInfoHandler contains all the requirements of the handler
type UpdateUserInfoHandler struct {
	Log   *zap.SugaredLogger
	Users user.CoreAPI
}

func (h UpdateUserInfoHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	userId := request.QueryStringParameters["id"]
	var reqUser user.User

	err := json.Unmarshal([]byte(request.Body), &reqUser)

	if err != nil {
		return events.APIGatewayProxyResponse{}, web.NewRequestError(err, http.StatusBadRequest)
	}

	dbUser, err := h.Users.Update(reqUser, userId)

	if err != nil {
		return events.APIGatewayProxyResponse{}, web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Respond(dbUser, http.StatusOK)
}
