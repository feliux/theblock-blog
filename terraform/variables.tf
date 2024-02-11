// Default
variable "custom_tags" {
  default = {
    app       = "theblock-blog"
    terraform = true
  }
}

// KMS
variable "kms" {
  description = "KMS options"
  default = {
    name        = "theblock-blog"
    description = "KMS for theblock-blog app"
    grant_name  = "lambda-theblock-blog"
  }
}

// SecretsManager
variable "sm" {
  description = "Name of the secrets to use by lambda"
  default = {
    description = "Secrets for theblock-blog app"
  }
}

// IAM
variable "iam" {
  description = "Options for IAM"
  default = {
    lambda_iam_role_name        = "theblock-blog"
    lambda_iam_policy_apgw      = "theblock-blog-apgw"
    lambda_iam_policy_sm        = "theblock-blog-secretsmanager"
    lambda_iam_policy_kms       = "theblock-blog-kms"
    lambda_iam_policy_s3        = "theblock-blog-s3"
    lambda_iam_policy_log_group = "theblock-blog-AWSLambdaBasicExecutionRole"
  }
}

// ApiGateway
variable "apgw" {
  description = "Options for APGW"
  default = {
    name                = "theblock-blog"
    description         = "theblock-blog app endpoints"
    path_part           = "{blog+}"
    stage_name          = "blog" // same as URL_PREFIX lambda env
    retention_logs_days = 7
  }
}

// Lambda
variable "lambda_conf" {
  description = "AWS Lambdas configuration"
  default = {
    lambda_name        = "theblock-blog"
    lambda_description = "Backend for theblock-blog app"
    //lambda_source_code  = "src/golang/get_recent_trades/main.go"
    //lambda_code_to_zip  = "src/golang/get_recent_trades/bootstrap"
    lambda_code_zipped = "zip/lambda.zip"
    lambda_handler     = "bootstrap"
    lambda_timeout     = 15  # seconds
    lambda_memory      = 128 # MB
    lambda_runtime     = "go1.x"
    //go_working_dir      = "src/golang/get_recent_trades" // for command go build
    //scheduled_rule      = true
    retention_logs_days = 7
    lambda_environment = {
      SECRET_NAME = "theblock-blog/app"    // also for sm name
      BUCKET_NAME = "theblock-blog" // also for bucket name
      URL_PREFIX  = "blog" // check on main.go method request path: {blog=login}
    }
  }
}
