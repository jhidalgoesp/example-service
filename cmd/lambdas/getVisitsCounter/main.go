// Package getVisitsCounter provides the entrypoint the service.
package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jhidalgoesp/example-services/cmd/lambdas/getVisitsCounter/handler"
	"github.com/jhidalgoesp/example-services/internal/awsProvider"
	"github.com/jhidalgoesp/example-services/internal/logger"
	"github.com/jhidalgoesp/example-services/pkg/counter"
	"github.com/jhidalgoesp/example-services/pkg/middleware"
	"os"
)

const tag = "getVisitsCounter"

func main() {
	log, err := logger.InitLogger(tag)
	if err != nil {
		fmt.Println("Error constructing logger:", err)
		os.Exit(1)
	}
	defer log.Sync()

	provider := awsProvider.NewProvider()
	dynamoClient := provider.NewDynamoClient()

	counter := counter.NewCore(log, dynamoClient)

	lambda.Start(
		middleware.Logger(
			log, middleware.Errors(
				log, handler.GetVisitsCounterHandler{Log: log, Core: counter}.Handle)))
}
