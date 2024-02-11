package routers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/feliux/theblock-blog/database"
	"github.com/feliux/theblock-blog/models"

	"github.com/aws/aws-lambda-go/events"
)

func NewRelation(ctx context.Context, request events.APIGatewayProxyRequest, claim models.Claim) models.Response {
	var response models.Response
	response.StatusCode = 400

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Message = "Parameter 'id' is mandatory."
		return response
	}

	var relation models.Relationship
	relation.UserId = claim.Id.Hex()
	relation.UserRelationId = id

	ok, err := database.InsertRelation(relation)
	if err != nil {
		response.Message = "Fail inserting relationship."
		log.Printf("Failed inserting relationship on database with ERROR: %s", err.Error())
		return response
	}

	if !ok {
		response.Message = "Fail inserting relationship."
		log.Println("Failed inserting relationship...")
		return response
	}

	response.StatusCode = 200
	response.Message = "Relation uploaded sucessfully."
	return response
}

func DeleteRelation(request events.APIGatewayProxyRequest, claim models.Claim) models.Response {
	var response models.Response
	response.StatusCode = 400

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Message = "Parameter 'id' is mandatory."
		return response
	}

	var relation models.Relationship
	relation.UserId = claim.Id.Hex()
	relation.UserRelationId = id

	ok, err := database.DeleteRelation(relation)
	if err != nil {
		response.Message = "Fail deleting relationship."
		log.Printf("Failed deleting relationship with ERROR: %s", err.Error())
		return response
	}
	if !ok {
		response.Message = "Fail deleting relationship."
		log.Println("Failed deleting relationship...")
		return response
	}

	response.StatusCode = 200
	response.Message = "Relation deleted sucessfully."
	return response

}

func GetRelation(request events.APIGatewayProxyRequest, claim models.Claim) models.Response {
	var response models.Response
	response.StatusCode = 400

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Message = "Parameter 'id' is mandatory."
		return response
	}

	var relation models.Relationship
	relation.UserId = claim.Id.Hex()
	relation.UserRelationId = id

	var resp models.ResponseRelationship

	existRelation := database.GetRelation(relation)
	if !existRelation {
		resp.Status = false
	} else {
		resp.Status = true
	}

	respJson, err := json.Marshal(existRelation)
	if err != nil {
		response.StatusCode = 500
		response.Message = "Failed getting data as json."
		log.Printf("Failed marshaling data into json format with ERROR: %s", err.Error())
		return response
	}

	response.StatusCode = 200
	response.Message = string(respJson)
	return response
}
