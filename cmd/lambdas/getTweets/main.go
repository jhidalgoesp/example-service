// Package getTweets provides the entrypoint the service.
package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/jhidalgoesp/example-services/cmd/lambdas/getTweets/handler"
	"github.com/jhidalgoesp/example-services/internal/awsProvider"
	"github.com/jhidalgoesp/example-services/internal/logger"
	"github.com/jhidalgoesp/example-services/pkg/middleware"
	"github.com/jhidalgoesp/example-services/pkg/twitter"
	"go.uber.org/zap"
	"os"
)

const tag = "getTweets"

func main() {
	log, err := logger.InitLogger(tag)
	if err != nil {
		fmt.Println("Error constructing logger:", err)
		os.Exit(1)
	}
	defer log.Sync()

	provider := awsProvider.NewProvider()
	SSMClient := provider.NewSSMClient()

	loadEnv(log, SSMClient)

	client := twitter.Client{Log: log}

	lambda.Start(
		middleware.Logger(
			log, middleware.Errors(
				log, handler.LastTweetsHandler{Log: log, Client: client}.Handle)))
}

func loadEnv(log *zap.SugaredLogger, SSMClient *ssm.SSM) {
	keys := []string{"consumerKey", "consumerSecret", "token", "tokenSecret"}

	for _, key := range keys {
		result, err := SSMClient.GetParameter(&ssm.GetParameterInput{
			Name: aws.String(fmt.Sprintf("/test-services/%s", key)),
		})

		if err != nil {
			log.Error(err.Error())
		}

		os.Setenv(key, *result.Parameter.Value)
	}
}
