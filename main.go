package main

import (
	"context"
	"os"
	"strings"

	"github.com/feliux/theblock-blog/awsauth"
	"github.com/feliux/theblock-blog/database"
	"github.com/feliux/theblock-blog/handlers"
	"github.com/feliux/theblock-blog/models"
	"github.com/feliux/theblock-blog/sm"

	"github.com/aws/aws-lambda-go/events"
	lambda "github.com/aws/aws-lambda-go/lambda"
)

var (
	region string = "us-east-1"
)

func ValidateParams() bool {
	var existsParam bool
	_, existsParam = os.LookupEnv("SECRET_NAME")
	if !existsParam {
		return existsParam
	}
	_, existsParam = os.LookupEnv("BUCKET_NAME")
	if !existsParam {
		return existsParam
	}
	_, existsParam = os.LookupEnv("URL_PREFIX")
	if !existsParam {
		return existsParam
	}
	return existsParam
}

func ExecuteLambda(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	var res *events.APIGatewayProxyResponse
	awsauth.InitAWS()
	if !ValidateParams() {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error loading env variables. It must include 'SECRET_NAME', 'BUCKET_NAME', 'URL_PREFIX'",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}
	secretModel, err := sm.GetSecret(os.Getenv("SECRET_NAME"))
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Error reading secret: " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}
	// Method request path: {blog=login}
	// Carefull when deploy with terraform. Check lambda env vars on variables.tf
	path := strings.Replace(request.PathParameters["blog"], os.Getenv("URL_PREFIX"), "", -1)
	awsauth.Ctx = context.WithValue(awsauth.Ctx, models.Key("path"), path)
	awsauth.Ctx = context.WithValue(awsauth.Ctx, models.Key("method"), request.HTTPMethod)
	awsauth.Ctx = context.WithValue(awsauth.Ctx, models.Key("body"), request.Body)
	awsauth.Ctx = context.WithValue(awsauth.Ctx, models.Key("bucketName"), os.Getenv("BUCKET_NAME"))
	awsauth.Ctx = context.WithValue(awsauth.Ctx, models.Key("host"), secretModel.Host)
	awsauth.Ctx = context.WithValue(awsauth.Ctx, models.Key("user"), secretModel.UserName)
	awsauth.Ctx = context.WithValue(awsauth.Ctx, models.Key("password"), secretModel.Password)
	awsauth.Ctx = context.WithValue(awsauth.Ctx, models.Key("jwtSign"), secretModel.JwtSign)
	awsauth.Ctx = context.WithValue(awsauth.Ctx, models.Key("database"), secretModel.Database)
	// transfer region to other packages
	awsauth.Region = region
	awsauth.Ctx = context.WithValue(awsauth.Ctx, models.Key("region"), region)

	err = database.Connect(awsauth.Ctx)
	if err != nil {
		res = &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error connecting to database: " + err.Error(),
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	}
	restAPI := handlers.Handlers(awsauth.Ctx, request)
	if restAPI.CustomResponse == nil { // si la respuesta no viene personalizada desde la api
		res = &events.APIGatewayProxyResponse{
			StatusCode: restAPI.StatusCode,
			Body:       restAPI.Message,
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
		}
		return res, nil
	} else {
		return restAPI.CustomResponse, nil
	}
}

func main() {
	// https://docs.aws.amazon.com/es_es/lambda/latest/dg/lambda-golang.html#golang-libraries
	lambda.Start(ExecuteLambda) // Llamada y comienzo de la lambda
}
