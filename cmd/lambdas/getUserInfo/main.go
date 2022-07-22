// Package getUserInfo provides the entrypoint the service.
package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/jhidalgoesp/example-services/cmd/lambdas/getUserInfo/handler"
	"github.com/jhidalgoesp/example-services/internal/awsProvider"
	"github.com/jhidalgoesp/example-services/internal/logger"
	"github.com/jhidalgoesp/example-services/pkg/counter"
	"github.com/jhidalgoesp/example-services/pkg/middleware"
	"github.com/jhidalgoesp/example-services/pkg/user"
	"os"
)

const tag = "getUserInfo"

func main() {
	log, err := logger.InitLogger(tag)
	if err != nil {
		fmt.Println("Error constructing logger:", err)
		os.Exit(1)
	}
	defer log.Sync()
	dynamoClient := awsProvider.NewProvider().NewDynamoClient()

	users := user.NewCore(log, dynamoClient)
	counter := counter.NewCore(log, dynamoClient)

	lambda.Start(
		middleware.Logger(
			log, middleware.Errors(
				log, handler.GetUserInfoHandler{Log: log, Users: users, Counter: counter}.Handle)))
}
