package routers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/feliux/theblock-blog/database"
	"github.com/feliux/theblock-blog/jwt"
	"github.com/feliux/theblock-blog/models"

	"github.com/aws/aws-lambda-go/events"
)

func Login(ctx context.Context) models.Response {
	var user models.User
	var response models.Response
	response.StatusCode = 400
	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		response.Message = "Invalid user or password."
		log.Printf("Error proccesing endpoint '%s' with ERROR: %s", ctx.Value(models.Key("path")).(string), err.Error())
		return response
	}
	if len(user.Email) == 0 {
		response.Message = "You must specify an email."
		log.Printf("User '%s' did not specify an email...", user.Email)
		return response
	}
	userData, userDataExists := database.TryLogin(user.Email, user.Password)
	if !userDataExists {
		response.Message = "Invalid user or password."
		log.Printf("Data for user '%s' does not exists on database...", user.Email)
		return response
	}
	jwtKey, err := jwt.GenerateJWT(ctx, userData)
	if err != nil {
		response.Message = "Error generating token jwt."
		log.Printf("Could not generate token for user '%s' with ERROR: ", user.Email, err.Error())
		return response
	}
	responseLogin := models.Login{
		Token: jwtKey,
	}
	token, err := json.Marshal(responseLogin)
	if err != nil {
		response.Message = "Error at marshall token."
		log.Printf("Could not marshall token to json format for user '%s' with ERROR: ", user.Email, err.Error())
		return response
	}
	cookie := &http.Cookie{
		Name:    "token",
		Value:   jwtKey,
		Expires: time.Now().Add(time.Hour * 24),
	}
	cookieStr := cookie.String()
	resp := &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(token),
		Headers: map[string]string{
			"Content-Type":                "application/json",
			"Access-Control-Allow-Origin": "*",
			"Set-Cookie":                  cookieStr,
		},
	}
	resp.StatusCode = 200
	response.Message = string(token)
	response.CustomResponse = resp
	return response
}
