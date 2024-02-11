package routers

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/feliux/theblock-blog/database"
	"github.com/feliux/theblock-blog/models"

	"github.com/aws/aws-lambda-go/events"
)

func GetUsers(request events.APIGatewayProxyRequest, claim models.Claim) models.Response {
	var response models.Response
	response.StatusCode = 400

	page := request.QueryStringParameters["page"]
	userType := request.QueryStringParameters["type"]
	search := request.QueryStringParameters["search"]
	userId := claim.Id.Hex()

	if len(page) < 1 {
		page = "1"
	}
	pag, err := strconv.Atoi(page)
	if err != nil {
		response.Message = "Parameter 'page' must be higher than 0."
		return response
	}
	allUsers, ok := database.GetAllUsers(userId, int64(pag), search, userType)
	if !ok {
		response.Message = "Error reading users."
		log.Println("Failed reading all users...")
		return response
	}
	respJson, err := json.Marshal(allUsers)
	if err != nil {
		response.StatusCode = 500
		response.Message = "Failed formating users data into json."
		log.Printf("Failed formating users data into json with ERROR: %s", err.Error())
		return response
	}
	response.StatusCode = 200
	response.Message = string(respJson)
	return response
}
