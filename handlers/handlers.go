package handlers

import (
	"context"
	"log"

	"github.com/feliux/theblock-blog/jwt"
	"github.com/feliux/theblock-blog/models"
	"github.com/feliux/theblock-blog/routers"

	"github.com/aws/aws-lambda-go/events"
)

func Handlers(ctx context.Context, request events.APIGatewayProxyRequest) models.Response {
	var response models.Response
	response.StatusCode = 400
	log.Println("Processing " + ctx.Value(models.Key("method")).(string) + " > " + ctx.Value(models.Key("path")).(string))
	//isOk, statusCode, msg, claim := validateAuthorization(ctx, request)
	isOk, statusCode, msg, claim := validateAuthorization(ctx, request)
	if !isOk {
		response.StatusCode = statusCode
		response.Message = msg
		return response
	}
	path := ctx.Value(models.Key("path")).(string)
	switch ctx.Value(models.Key("method")).(string) {
	case "GET":
		switch path {
		case "profile":
			return routers.ViewProfile(request)
		case "tweet":
			return routers.ReadTweet(request)
		case "avatar":
			return routers.GetImage(ctx, "A", request, claim)
		case "banner":
			return routers.GetImage(ctx, "B", request, claim)
		case "relation":
			return routers.GetRelation(request, claim)
		case "users":
			return routers.GetUsers(request, claim)
		case "followersTweets":
			return routers.GetFollowerTweets(request, claim)
		}
	case "POST":
		switch path {
		case "register":
			return routers.RegisterUser(ctx)
		case "login":
			return routers.Login(ctx)
		case "tweet":
			return routers.SendTweet(ctx, claim)
		case "avatar":
			return routers.UploadImage(ctx, "A", request, claim)
		case "banner":
			return routers.UploadImage(ctx, "B", request, claim)
		case "relation":
			return routers.NewRelation(ctx, request, claim)
		}
	case "PUT":
		switch path {
		case "modify":
			return routers.ModifyProfile(ctx, claim)
		}
	case "DELETE":
		switch path {
		case "tweet":
			return routers.DeleteTweet(request, claim)
		case "relation":
			return routers.DeleteRelation(request, claim)
		}
	}
	response.Message = "Invalid method."
	return response
}

func validateAuthorization(ctx context.Context, request events.APIGatewayProxyRequest) (bool, int, string, models.Claim) {
	path := ctx.Value(models.Key("path")).(string)
	if path == "register" || path == "login" || path == "getAvatar" || path == "getBanner" {
		return true, 200, "", models.Claim{}
	}
	token := request.Headers["Authorization"]
	if len(token) == 0 {
		return false, 401, "Token is required.", models.Claim{}
	}

	claim, todoOK, msg, err := jwt.ProccesToken(token, ctx.Value(models.Key("jwtSign")).(string))
	if !todoOK {
		if err != nil {
			log.Println("Failed processing token with ERROR: " + err.Error())
			return false, 401, err.Error(), models.Claim{}
		} else {
			log.Println("Failed processing token wit ERROR: " + msg)
			return false, 401, msg, models.Claim{}
		}
	}
	log.Println("Token validation is correct...")
	return true, 200, msg, *claim
}
