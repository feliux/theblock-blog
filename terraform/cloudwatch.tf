// Cloudwatch logs
resource "aws_cloudwatch_log_group" "lambda_logs" {
  name              = "/aws/lambda/${lookup(var.lambda_conf, "lambda_name")}"
  retention_in_days = lookup(var.lambda_conf, "retention_logs_days")
}

resource "aws_cloudwatch_log_group" "apgw_logs" {
  // the aws_cloudwatch_log_group resource can be used where the name matches the API Gateway naming convention
  name              = "API-Gateway-Execution-Logs_${aws_api_gateway_rest_api.this.id}/${var.apgw.stage_name}"
  retention_in_days = lookup(var.apgw, "retention_logs_days")
}
