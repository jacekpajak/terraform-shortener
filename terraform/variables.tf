variable "aws_region" {
  description = "The AWS region where resources will be created."
  default     = "eu-central-1"
}

variable "lambda_function_name" {
  description = "The name of the Lambda function."
  default     = "url-shortener"
}

variable "dynamodb_table_name" {
  description = "The name of the DynamoDB table."
  default     = "urls"
}

variable "s3_bucket_name" {
  description = "The S3 bucket for Lambda deployment package."
  default     = "lambda_deployment_package"
}