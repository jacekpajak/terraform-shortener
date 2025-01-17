output "lambda_invoke_url" {
  description = "Invoke URL for the Lambda function via API Gateway."
  value       = aws_lambda_function.url_shortener.invoke_arn
}

output "dynamodb_table_name" {
  description = "The name of the DynamoDB table."
  value       = aws_dynamodb_table.urls.name
}

output "api_gateway_url" {
  value = aws_api_gateway_deployment.deployment.invoke_url
}