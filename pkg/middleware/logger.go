package middleware

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/jhidalgoesp/example-services/pkg/web"
	"go.uber.org/zap"
)

// Logger handles logging requests and responses coming  in and out of the call chain.
func Logger(log *zap.SugaredLogger, handler web.Handler) web.Handler {
	return func(r events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		log.Infow("request started", "method", r.HTTPMethod, "path", r.Path,
			"remoteaddr", r.RequestContext.Identity.SourceIP, "body", r.Body, "queryParams", r.QueryStringParameters)
		response, err := handler(r)
		log.Infow("request completed", "method", r.HTTPMethod, "path", r.Path,
			"remoteaddr", r.RequestContext.Identity.SourceIP, "response", response)
		return response, err
	}
}
