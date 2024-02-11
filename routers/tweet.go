package routers

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/feliux/theblock-blog/database"
	"github.com/feliux/theblock-blog/models"

	"github.com/aws/aws-lambda-go/events"
)

func SendTweet(ctx context.Context, claim models.Claim) models.Response {
	var tweet models.Tweet
	var response models.Response
	response.StatusCode = 400
	userId := claim.Id.Hex()
	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &tweet)
	if err != nil {
		response.Message = "Incorrect data body."
		log.Printf("Error proccesing endpoint '%s' with ERROR: %s", ctx.Value(models.Key("path")).(string), err.Error())
		return response
	}
	register := models.SaveTweet{
		UserId:  userId,
		Message: tweet.Message,
		Date:    time.Now(),
	}
	_, status, err := database.InsertTweet(register)
	if err != nil {
		response.Message = "Insert tweet operation failed."
		log.Printf("Insert tweet operation failed with ERROR: %s", err.Error())
		return response
	}
	if !status {
		response.Message = "Not able to insert tweet."
		log.Printf("Not able to insert tweet on database")
		return response
	}
	response.StatusCode = 200
	response.Message = "Tweet created succesfully."
	return response
}

func ReadTweet(request events.APIGatewayProxyRequest) models.Response {
	var response models.Response
	response.StatusCode = 400
	id := request.QueryStringParameters["id"]
	page := request.QueryStringParameters["page"]
	if len(id) < 1 {
		response.Message = "Parameter 'id' is mandatory."
		return response
	}
	if len(page) < 1 {
		page = "1"
	}
	pag, err := strconv.Atoi(page)
	if err != nil {
		response.Message = "Parameter 'page' must be higher than 0."
		return response
	}
	tweets, ok := database.ReadTweets(id, int64(pag))
	if !ok {
		response.Message = "Error reading tweets."
		log.Println("Failed reading tweets...")
		return response
	}
	respJson, err := json.Marshal(tweets)
	if err != nil {
		response.StatusCode = 500
		response.Message = "Failed formating tweets data into json."
		log.Printf("Failed formating tweets data into json with ERROR: %s", err.Error())
		return response
	}
	response.StatusCode = 200
	response.Message = string(respJson)
	return response
}

func DeleteTweet(request events.APIGatewayProxyRequest, claim models.Claim) models.Response {
	var response models.Response
	response.StatusCode = 400
	tweetId := request.QueryStringParameters["id"]
	if len(tweetId) < 1 {
		response.Message = "Parameter 'id' is mandatory."
		return response
	}
	err := database.DeleteTweet(tweetId, claim.Id.Hex())
	if err != nil {
		response.Message = "Failed deleting tweet."
		log.Printf("Failed deleting tweets with ERROR: %s", err.Error())
		return response
	}
	response.Message = "Tweet deleted sucessfully."
	response.StatusCode = 200
	return response
}
