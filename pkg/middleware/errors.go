package middleware

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/jhidalgoesp/example-services/pkg/validate"
	"github.com/jhidalgoesp/example-services/pkg/web"
	"go.uber.org/zap"
	"net/http"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.
func Errors(log *zap.SugaredLogger, handler web.Handler) web.Handler {
	h := func(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		response, err := handler(r)
		if err != nil {
			log.Errorw("ERROR", err.Error())
			var er web.ErrorResponse
			var status int
			switch act := validate.Cause(err).(type) {
			case *web.RequestError:
				er = web.ErrorResponse{
					Error: act.Error(),
				}
				status = act.Status
			default:
				er = web.ErrorResponse{
					Error: http.StatusText(http.StatusInternalServerError),
				}
				status = http.StatusInternalServerError
			}
			return web.Respond(er, status)
		}
		return response, nil
	}
	return h
}
