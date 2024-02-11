package models

import (
	"github.com/aws/aws-lambda-go/events"
)

type Response struct {
	StatusCode     int
	Message        string
	CustomResponse *events.APIGatewayProxyResponse
}
