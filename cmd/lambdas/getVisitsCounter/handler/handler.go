// Package handler provides the http logic of the service.
package handler

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/jhidalgoesp/example-services/pkg/counter"
	"github.com/jhidalgoesp/example-services/pkg/web"
	"go.uber.org/zap"
	"net/http"
)

// GetVisitsCounterHandler contains all the requirements of the handler
type GetVisitsCounterHandler struct {
	Log  *zap.SugaredLogger
	Core counter.CoreAPI
}

func (h GetVisitsCounterHandler) Handle(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	counter := counter.Counter{}
	err := h.Core.GetVisitsCounter(&counter)

	if err != nil {
		return events.APIGatewayProxyResponse{}, web.NewRequestError(err, http.StatusBadRequest)
	}

	return web.Respond(counter, http.StatusOK)
}
