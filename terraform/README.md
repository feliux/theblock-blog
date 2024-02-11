Create the following files in order to deploy with the current configuration

**./secrets/creds**
~~~
[develop]
aws_access_key_id = changeme
aws_secret_access_key = changeme
~~~

**secrets.tf**
```tf
variable "credentials" {
  description = "Service account credentials"
  default     = "./secrets/creds"
}

variable "profile" {
  default = "develop"
}

variable "accountId" {
  description = "Account for cloud resources"
  default     = "changeme" // your aws account id
}

variable "region" {
  description = "Region for cloud resources"
  default     = "us-east-1" // region to deploy
}

variable "secrets" {
  type = map(string)
  default = {
    host     = "changeme" // <yourMongoHost: mydb.xxxxxx.mongodb.net>
    database = "changeme" // <yourMongoUser: root>
    username = "changeme" // <yourMongoPassword: changeme>
    password = "changeme" // <yourMongoDatabase: mydb>
    jwtSign  = "changeme" // <strongPasswor: wiuebfwuiebf564we89f>
  }
}
```
