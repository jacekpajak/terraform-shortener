package handlers

import (
	"url-shortener/database" // Import the database package

	"github.com/aws/aws-lambda-go/events"
)

// RedirectURLHandlerAPI handles redirection for Lambda-based API Gateway
func RedirectURLHandlerAPI(db database.Database, shortURL string) (events.APIGatewayProxyResponse, error) {
	// Validate short URL
	if shortURL == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Short URL is required",
		}, nil
	}

	// Retrieve the original URL from the database
	originalURL, err := db.GetURL(shortURL)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
			Body:       "URL not found",
		}, nil
	}

	// Redirect to the original URL
	return events.APIGatewayProxyResponse{
		StatusCode: 301,
		Headers: map[string]string{
			"Location": originalURL,
		},
	}, nil
}
