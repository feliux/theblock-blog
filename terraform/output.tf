output "apigw_invoke_url" {
  value = aws_api_gateway_stage.this.invoke_url
}

output "aws_lambda_permission_source_arn" {
  value = "arn:aws:execute-api:${var.region}:${var.accountId}:${aws_api_gateway_rest_api.this.id}/*/*${aws_api_gateway_resource.this.path}"
}
