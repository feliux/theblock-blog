resource "aws_api_gateway_rest_api" "this" {
  name        = var.apgw.name
  description = var.apgw.description
  endpoint_configuration {
    types = ["REGIONAL"]
  }
  binary_media_types = ["multipart/form-data", "image/jpeg"]
  tags               = var.custom_tags
}

resource "aws_api_gateway_resource" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  parent_id   = aws_api_gateway_rest_api.this.root_resource_id
  path_part   = var.apgw.path_part
}

resource "aws_api_gateway_method" "any" {
  rest_api_id   = aws_api_gateway_rest_api.this.id
  resource_id   = aws_api_gateway_resource.this.id
  http_method   = "ANY"
  authorization = "NONE"
}

resource "aws_api_gateway_integration" "this" {
  rest_api_id             = aws_api_gateway_rest_api.this.id
  resource_id             = aws_api_gateway_resource.this.id
  http_method             = aws_api_gateway_method.any.http_method
  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.lambdas.invoke_arn
}

resource "aws_api_gateway_deployment" "this" {
  rest_api_id = aws_api_gateway_rest_api.this.id
  description = "New deployment for api"
  triggers = {
    redeployment = sha1(jsonencode([
      aws_api_gateway_resource.this.id,
      aws_api_gateway_method.any.id,
      aws_api_gateway_integration.this.id,
      filebase64sha256("${path.module}/${lookup(var.lambda_conf, "lambda_code_zipped")}"),
    ]))
  }
  lifecycle {
    create_before_destroy = true
  }
}

resource "aws_api_gateway_stage" "this" {
  deployment_id = aws_api_gateway_deployment.this.id
  rest_api_id   = aws_api_gateway_rest_api.this.id
  stage_name    = var.apgw.stage_name
  description   = var.apgw.description
  tags          = var.custom_tags
  depends_on    = [aws_cloudwatch_log_group.apgw_logs]
}