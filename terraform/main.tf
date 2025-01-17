# DynamoDB Table
resource "aws_dynamodb_table" "urls" {
  name         = var.dynamodb_table_name
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "short_url"

  attribute {
    name = "short_url"
    type = "S"
  }
}

# IAM Role for Lambda
resource "aws_iam_role" "lambda_exec_role" {
  name = "lambda_exec_role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "lambda.amazonaws.com"
        },
        Action = "sts:AssumeRole"
      }
    ]
  })
}

# Policy for DynamoDB Access and CloudWatch Logs
resource "aws_iam_policy" "lambda_policy" {
  name = "lambda_policy"
  policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Action = [
          "dynamodb:PutItem",
          "dynamodb:GetItem",
          "dynamodb:Query"
        ],
        Resource = aws_dynamodb_table.urls.arn
      },
      {
        Effect = "Allow",
        Action = [
          "logs:CreateLogGroup",
          "logs:CreateLogStream",
          "logs:PutLogEvents"
        ],
        Resource = "arn:aws:logs:*:*:log-group:/aws/lambda/*"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "lambda_attach" {
  role       = aws_iam_role.lambda_exec_role.name
  policy_arn = aws_iam_policy.lambda_policy.arn
}

# Lambda Function
resource "aws_lambda_function" "url_shortener" {
  function_name = var.lambda_function_name
  runtime       = "provided.al2"
  handler       = "bootstrap"
  role          = aws_iam_role.lambda_exec_role.arn

  filename         = "../deploy.zip"
  source_code_hash = filebase64sha256("../deploy.zip")

  environment {
    variables = {
      TABLE_NAME = var.dynamodb_table_name
      AWS_EXECUTION_ENV = "1"
    }
  }
}

# API Gateway
resource "aws_api_gateway_rest_api" "url_shortener_api" {
  name        = "url-shortener-api"
  description = "API Gateway for URL shortener service"

  endpoint_configuration {
    types = ["REGIONAL"]
  }
}

# Enable CloudWatch Logging
resource "aws_api_gateway_account" "api_logging" {
  cloudwatch_role_arn = aws_iam_role.api_gateway_logging_role.arn
}

resource "aws_iam_role" "api_gateway_logging_role" {
  name = "api-gateway-logging-role"

  assume_role_policy = jsonencode({
    Version = "2012-10-17",
    Statement = [
      {
        Effect = "Allow",
        Principal = {
          Service = "apigateway.amazonaws.com"
        },
        Action = "sts:AssumeRole"
      }
    ]
  })
}

resource "aws_iam_role_policy_attachment" "api_gateway_logging_policy_attach" {
  role       = aws_iam_role.api_gateway_logging_role.name
  policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonAPIGatewayPushToCloudWatchLogs"
}

# POST Method (Shorten URL)
resource "aws_api_gateway_resource" "shorten" {
  rest_api_id = aws_api_gateway_rest_api.url_shortener_api.id
  parent_id   = aws_api_gateway_rest_api.url_shortener_api.root_resource_id
  path_part   = "shorten"
}

resource "aws_api_gateway_method" "post_shorten" {
  rest_api_id   = aws_api_gateway_rest_api.url_shortener_api.id
  resource_id   = aws_api_gateway_resource.shorten.id
  http_method   = "POST"
  authorization = "AWS_IAM"
}

resource "aws_api_gateway_method_settings" "post_logging" {
  rest_api_id = aws_api_gateway_rest_api.url_shortener_api.id
  stage_name  = aws_api_gateway_stage.prod.stage_name
  method_path = "shorten/POST"

  settings {
    logging_level      = "INFO"
    metrics_enabled    = true
    data_trace_enabled = true
  }
}

resource "aws_api_gateway_integration" "post_shorten_lambda" {
  rest_api_id = aws_api_gateway_rest_api.url_shortener_api.id
  resource_id = aws_api_gateway_resource.shorten.id
  http_method = aws_api_gateway_method.post_shorten.http_method

  integration_http_method = "POST"
  type                    = "AWS_PROXY"
  uri                     = aws_lambda_function.url_shortener.invoke_arn
}

# Deployment
resource "aws_api_gateway_deployment" "deployment" {
  rest_api_id = aws_api_gateway_rest_api.url_shortener_api.id

  triggers = {
    redeployment = sha1(jsonencode(aws_api_gateway_rest_api.url_shortener_api.body))
  }
}

resource "aws_api_gateway_stage" "prod" {
  stage_name    = "prod"
  rest_api_id   = aws_api_gateway_rest_api.url_shortener_api.id
  deployment_id = aws_api_gateway_deployment.deployment.id

  depends_on = [aws_api_gateway_account.api_logging]

  access_log_settings {
    destination_arn = aws_cloudwatch_log_group.api_gw_logs.arn
    format          = "$context.requestId $context.httpMethod $context.resourcePath $context.status $context.integrationErrorMessage"
  }
}

resource "aws_cloudwatch_log_group" "api_gw_logs" {
  name              = "/aws/api-gateway/url-shortener-api"
  retention_in_days = 7
}

data "aws_caller_identity" "current" {}

resource "aws_lambda_permission" "api_gateway_invoke" {
  statement_id  = "AllowApiGatewayInvoke"
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.url_shortener.function_name
  principal     = "apigateway.amazonaws.com"
  source_arn    = "arn:aws:execute-api:${var.aws_region}:${data.aws_caller_identity.current.account_id}:${aws_api_gateway_rest_api.url_shortener_api.id}/*/*"
}