package sm

import (
	"encoding/json"
	"log"

	"github.com/feliux/theblock-blog/awsauth"
	"github.com/feliux/theblock-blog/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetSecret(secretName string) (models.Secret, error) {
	var secretData models.Secret
	svc := secretsmanager.NewFromConfig(awsauth.Cfg)
	key, err := svc.GetSecretValue(awsauth.Ctx, &secretsmanager.GetSecretValueInput{
		SecretId: aws.String(secretName),
	})
	if err != nil {
		panic("Error getting secrets... > " + err.Error())
		return secretData, err
	}
	json.Unmarshal([]byte(*key.SecretString), &secretData)
	log.Println("Secrets retrieved succesfully...")
	return secretData, nil
}
