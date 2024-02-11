package routers

import (
	"context"
	"encoding/json"
	"log"

	"github.com/feliux/theblock-blog/database"
	"github.com/feliux/theblock-blog/models"
)

func RegisterUser(ctx context.Context) models.Response {
	var user models.User
	var response models.Response
	response.StatusCode = 400
	body := ctx.Value(models.Key("body")).(string)
	err := json.Unmarshal([]byte(body), &user)
	if err != nil {
		response.Message = "Error processing body."
		log.Printf("Error processing endpoint '%s' with ERROR: %s", ctx.Value(models.Key("path")).(string), err.Error())
		return response
	}
	if len(user.Email) == 0 {
		response.Message = "You must specify an email."
		log.Println("User did not specify an email...")
		return response
	}
	if len(user.Password) < 8 {
		response.Message = "You must specify a password with more than 8 characters."
		log.Printf("User '%s' must specify a password with more than 8 characters...", user.Email)
		return response
	}
	_, userExists, _ := database.CheckUser(user.Email)
	if userExists {
		response.Message = "User already exists."
		log.Printf("User '%s' already exists on database...", user.Email)
		return response
	}
	_, commandStatus, err := database.Insert(user)
	if err != nil {
		response.Message = "Error registering user on database."
		log.Printf("Error registering user '%s' on database with ERROR: %s", user.Email, err.Error())
		return response
	}
	if !commandStatus {
		response.Message = "Error registering user on database."
		log.Printf("Error registering user '%s' on database with ERROR: general", user.Email)
		return response
	}
	response.StatusCode = 200
	response.Message = "User registered on database."
	log.Printf("User '%s' registered on database...", user.Email)
	return response
}
