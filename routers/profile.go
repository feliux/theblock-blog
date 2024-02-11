package routers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/feliux/theblock-blog/database"
	"github.com/feliux/theblock-blog/models"

	"github.com/aws/aws-lambda-go/events"
)

func ViewProfile(request events.APIGatewayProxyRequest) models.Response {
	var response models.Response
	response.StatusCode = 400
	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Message = "Parameter 'id' is mandatory."
		return response
	}
	profile, err := database.SearchProfile(id)
	if err != nil {
		response.Message = "Error fetching user profile."
		log.Printf("Error fetching user profile with ERROR: %s", err.Error())
		return response
	}
	respJson, err := json.Marshal(profile)
	if err != nil {
		response.StatusCode = 500
		response.Message = "Error marshaling user data profile."
		log.Printf("Error marshaling user data profile with ERROR: %s", err.Error())
		return response
	}
	response.StatusCode = 200
	response.Message = string(respJson)
	return response
}

func ModifyProfile(ctx context.Context, claim models.Claim) models.Response {
	var user models.User
	var response models.Response
	response.StatusCode = 400
	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		response.Message = "Incorrect data body."
		log.Printf("Error proccesing endpoint '%s' with ERROR: %s", ctx.Value(models.Key("path")).(string), err.Error())
		return response
	}
	status, err := database.ModifyRegister(user, claim.Id.Hex())
	if err != nil {
		response.Message = "Modify register operation failed."
		log.Printf("Error modifying register with ERROR: %s", err.Error())
		return response
	}
	if !status {
		response.Message = "Not able to modify user register."
		log.Printf("Not able to modify register on database")
		return response
	}
	response.StatusCode = 200
	response.Message = "Data updated succesfully."
	return response
}
