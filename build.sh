
rm bootstrap terraform/zip/lambda.zip 2>/dev/null
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o bootstrap main.go
zip terraform/zip/lambda.zip bootstrap # bootstrap binary is the AWS handler name
