package routers

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/feliux/theblock-blog/awsauth"
	"github.com/feliux/theblock-blog/database"
	"github.com/feliux/theblock-blog/models"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
)

func GetImage(ctx context.Context, fileType string, request events.APIGatewayProxyRequest, claim models.Claim) models.Response {

	var response models.Response
	response.StatusCode = 400

	id := request.QueryStringParameters["id"]
	if len(id) < 1 {
		response.Message = "Parameter 'id' is mandatory."
		return response
	}

	user, err := database.SearchProfile(id)
	if err != nil {
		response.Message = "User id does not exists."
		log.Printf("User does not exists with ERROR: %s", err.Error())
		return response
	}

	var filename string
	switch fileType {
	case "A":
		filename = user.Avatar
	case "B":
		filename = user.Banner
	}
	log.Printf("Getting filename: %s", filename)
	svc := s3.NewFromConfig(awsauth.Cfg)
	file, err := downloadFromS3(ctx, svc, filename)
	if err != nil {
		response.StatusCode = 500
		response.Message = "Fail downloading image."
		log.Printf("Fail downloading image with ERROR: %s", err.Error())
		return response
	}

	response.CustomResponse = &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       file.String(),
		Headers: map[string]string{
			"Content-Type":        "application/octet-stream",
			"Content-Disposition": fmt.Sprintf("attachment; filename=\"%s\"", filename),
		},
	}
	return response
}

func downloadFromS3(ctx context.Context, svc *s3.Client, filename string) (*bytes.Buffer, error) {
	bucket := ctx.Value(models.Key("bucketName")).(string)
	obj, err := svc.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(filename),
	})
	if err != nil {
		return nil, err
	}
	defer obj.Body.Close()
	log.Printf("Downloading image from bucket: %s", bucket)

	file, err := ioutil.ReadAll(obj.Body)
	if err != nil {
		return nil, err
	}
	buffer := bytes.NewBuffer(file)
	return buffer, nil
}
