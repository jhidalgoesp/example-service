package web

import (
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
)

// Respond converts a Go value to JSON wraps it into a APIGatewayResponse and sends it to the client.
func Respond(data interface{}, statusCode int) (events.APIGatewayProxyResponse, error) {
	// If there is nothing to marshal then set status code and return.
	if statusCode == http.StatusNoContent {
		return events.APIGatewayProxyResponse{
			Headers: map[string]string{
				"Content-Type":                 "application/json",
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "*",
			},
			StatusCode: statusCode,
		}, nil
	}

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Send the result back to the client.

	return events.APIGatewayProxyResponse{
		Body: string(jsonData),
		Headers: map[string]string{
			"Content-Type":                 "application/json",
			"Access-Control-Allow-Origin":  "*",
			"Access-Control-Allow-Methods": "*",
		},
		StatusCode: statusCode,
	}, nil
}
