# The Block

**DEPRECATED: infra is down and repo is public**

Blog service like Twitter built on golang.

- AWS Api Gateway to expose endpoints
- Lambda as golang backend
- KMS and SM for storing secrets
- MongoDB as database

## Deloy

First change the region where you want to deploy on [main.go](./main.go) file.

**MongoDB**

Create a Mongo server on [cloud.mongodb.com](https://cloud.mongodb.com/).

**Secrets Manager**

Your secret manager must include the following secrets. If deploy with terraform check the [terraform/README.md](./terraform/README.md)

~~~
host: <yourMongoHost: mydb.xxxxxx.mongodb.net>
username: <yourMongoUser: root>
password: <yourMongoPassword: changeme>
database: <yourMongoDatabase: mydb>
jwtSign: : <strongPasswor: wiuebfwuiebf564we89f>
~~~

**Terraform**

Just compile running `bash buld.sh`. It generates the `bootstrap` binary and the `terraform/zip/lambda.zip` for uploading on AWS Lambda. Then execute terraform.

```sh
# Build
$ bash build.sh
# Deploy
$ cd terraform
$ terraform init
$ terraform plan -out tfplan
$ terraform apply -auto-approve
$ terraform destroy -auto-approve
```

## Tests

Check [apiCalls](./apiCalls.md) for testing endpoints.
