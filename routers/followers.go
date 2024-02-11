package routers

import (
	"encoding/json"
	"log"
	"strconv"

	"github.com/feliux/theblock-blog/database"
	"github.com/feliux/theblock-blog/models"

	"github.com/aws/aws-lambda-go/events"
)

func GetFollowerTweets(request events.APIGatewayProxyRequest, claim models.Claim) models.Response {
	var response models.Response
	response.StatusCode = 400
	userId := claim.Id.Hex()
	page := request.QueryStringParameters["page"]
	if len(page) < 1 {
		page = "1"
	}
	pag, err := strconv.Atoi(page)
	if err != nil {
		response.Message = "Parameter 'page' must be higher than 0."
		return response
	}
	tweets, ok := database.GetFollowerTweets(userId, pag)
	if !ok {
		response.Message = "Error reading followers tweets."
		log.Println("Failed reading followers tweets...")
		return response
	}
	respJson, err := json.Marshal(tweets)
	if err != nil {
		response.StatusCode = 500
		response.Message = "Failed formating followers tweets data into json."
		log.Printf("Failed formating followers tweets data into json with ERROR: %s", err.Error())
		return response
	}
	response.StatusCode = 200
	response.Message = string(respJson)
	return response
}
