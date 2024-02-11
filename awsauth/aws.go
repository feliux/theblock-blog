package awsauth

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	Ctx    context.Context
	Cfg    aws.Config
	err    error
	Region string
)

func InitAWS() {
	Ctx = context.TODO()
	Cfg, err = config.LoadDefaultConfig(Ctx, config.WithDefaultRegion(Region))
	if err != nil {
		panic("Error loading default .aws/config: " + err.Error())
	}
}
