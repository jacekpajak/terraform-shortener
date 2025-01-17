# URL Shortener with Go and Terraform

## Overview
This repository contains a serverless **URL shortener service** built with **Go** and deployed using **AWS Lambda** and **API Gateway**. The infrastructure is managed with **Terraform**, and DynamoDB is used as the backend for storing URL mappings.

## Features
- **Shorten URLs**: Convert long URLs into short, easy-to-share links.
- **Redirect to Original URL**: Access the original URL using the short link.
- **Serverless Architecture**: Deployed on AWS Lambda and API Gateway for scalability.
- **Infrastructure as Code**: Fully automated deployment with Terraform.

## Endpoints
### POST `/shorten`
- **Description**: Accepts a long URL and returns a shortened URL.
- **Request Body**:
  ```json
  {
    "url": "https://example.com"
  }
