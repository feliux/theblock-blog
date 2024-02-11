resource "aws_lambda_function" "lambdas" {
  //count         = length(var.lambda_conf)
  filename      = "${path.module}/${lookup(var.lambda_conf, "lambda_code_zipped")}"
  package_type  = "Zip" // Image
  function_name = lookup(var.lambda_conf, "lambda_name")
  description   = lookup(var.lambda_conf, "lambda_description")
  role          = aws_iam_role.lambda_role.arn
  handler       = lookup(var.lambda_conf, "lambda_handler")
  timeout       = lookup(var.lambda_conf, "lambda_timeout")
  memory_size   = lookup(var.lambda_conf, "lambda_memory")

  # The filebase64sha256() function is available in Terraform 0.11.12 and later
  # For Terraform 0.11.11 and earlier, use the base64sha256() function and the file() function:
  # source_code_hash = "${base64sha256(file("lambda_function_payload.zip"))}"
  source_code_hash = filebase64sha256("${path.module}/${lookup(var.lambda_conf, "lambda_code_zipped")}")
  runtime = lookup(var.lambda_conf, "lambda_runtime")

  environment {
    variables = lookup(var.lambda_conf, "lambda_environment")
  }

  tags = var.custom_tags

  depends_on = [
    aws_iam_role_policy_attachment.attach_apgw_policy_to_iam_role,
    aws_iam_role_policy_attachment.attach_sm_policy_to_iam_role,
    aws_cloudwatch_log_group.lambda_logs
  ]
}

resource "aws_lambda_permission" "apigw_lambda" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.lambdas.function_name
  principal     = "apigateway.amazonaws.com"

  # More: http://docs.aws.amazon.com/apigateway/latest/developerguide/api-gateway-control-access-using-iam-policies-to-invoke-api.html
  //source_arn = "arn:aws:execute-api:${var.region}:${var.accountId}:${aws_api_gateway_rest_api.this.id}/*/${aws_api_gateway_method.any.http_method}${aws_api_gateway_resource.this.path}"
  source_arn = "arn:aws:execute-api:${var.region}:${var.accountId}:${aws_api_gateway_rest_api.this.id}/*/*/*"
}
