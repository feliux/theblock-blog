package routers

import (
	"bytes"
	"context"
	"encoding/base64"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"strings"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"

	"github.com/feliux/theblock-blog/database"
	"github.com/feliux/theblock-blog/models"
)

type readSeeker struct {
	io.Reader
}

func (rs *readSeeker) Seek(offset int64, whence int) (int64, error) {
	return 0, nil
}

func UploadImage(ctx context.Context, fileType string, request events.APIGatewayProxyRequest, claim models.Claim) models.Response {

	var response models.Response
	response.StatusCode = 400
	userId := claim.Id.Hex()

	var filename string
	var user models.User

	bucket := aws.String(ctx.Value(models.Key("bucketName")).(string))
	switch fileType {
	case "A":
		filename = "avatars/" + userId + ".jpg" // Rename with the userId
		user.Avatar = filename
	case "B":
		filename = "banners/" + userId + ".jpg" // Rename with the userId
		user.Banner = filename
	}

	mediaType, params, err := mime.ParseMediaType(request.Headers["content-type"])
	if err != nil {
		response.StatusCode = 500
		response.Message = "Can not parse 'Content-Type'."
		log.Printf("Can not parse 'Content-Type' with ERROR: %s", err.Error())
		return response
	}

	if strings.HasPrefix(mediaType, "multipart/") { // multipart (ie uploaded from form), not binary
		body, err := base64.StdEncoding.DecodeString(request.Body) // image come in the body
		if err != nil {
			response.StatusCode = 500
			response.Message = "Can not decode multipart."
			log.Printf("Can not decode multipart from request body with ERROR: %s", err.Error())
			return response
		}

		mr := multipart.NewReader(bytes.NewReader(body), params["boundary"])
		p, err := mr.NextPart()
		if err != nil && err != io.EOF {
			response.StatusCode = 500
			response.Message = "Failed on multipart."
			log.Printf("Failed on multipart with ERROR: %s", err.Error())
			return response
		}

		if err != io.EOF {
			if p.FileName() != "" {
				buf := bytes.NewBuffer(nil)
				if _, err := io.Copy(buf, p); err != nil {
					response.StatusCode = 500
					response.Message = "Failed converting image data."
					log.Printf("Failed converting image data with ERROR: %s", err.Error())
					return response
				}
				// Auth
				sess, err := session.NewSession(&aws.Config{
					// Region: aws.String(region)})
					Region: aws.String(ctx.Value(models.Key("region")).(string))})

				if err != nil {
					response.StatusCode = 500
					response.Message = "General error."
					log.Printf("Failed creating new aws session with ERROR: %s", err.Error())
					return response
				}

				uploader := s3manager.NewUploader(sess)
				_, err = uploader.Upload(&s3manager.UploadInput{
					Bucket: bucket,
					Key:    aws.String(filename),
					Body:   &readSeeker{buf},
				})

				if err != nil {
					response.StatusCode = 500
					response.Message = "Can not save image data."
					log.Printf("Failed uploading image to s3 with ERROR: %s", err.Error())
					return response
				}
			}
		}

		status, err := database.ModifyRegister(user, userId)
		if err != nil || !status {
			response.StatusCode = 400
			response.Message = "Modify register operation failed."
			log.Printf("Error modifying register with ERROR: %s", err.Error())
			return response
		}
	} else {
		response.StatusCode = 400
		response.Message = "Must send the image with the header 'Content-Type' of type 'multipart/'."
		log.Println("User did not send the image in the correct format...")
		return response
	}

	response.StatusCode = 200
	response.Message = "Image uploaded sucessfully."
	return response
}
