package handlers

import (
	"encoding/json"

	"url-shortener/database"
	"url-shortener/utils"

	"github.com/aws/aws-lambda-go/events"
)

// StoreURLHandlerAPI handles storing URLs in the database
func StoreURLHandlerAPI(db database.Database, body string) (events.APIGatewayProxyResponse, error) {
	// Parse JSON request body
	var input map[string]string
	err := json.Unmarshal([]byte(body), &input)
	if err != nil || input["url"] == "" {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       "Invalid input",
		}, nil
	}

	// Generate a short URL
	originalURL := input["url"]
	shortURL := utils.GenerateShortURL()

	// Store the URL in the database
	err = db.StoreURL(shortURL, originalURL)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Error storing URL",
		}, nil
	}

	// Prepare response
	response, err := json.Marshal(map[string]string{
		"short_url":    shortURL,
		"original_url": originalURL,
	})
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       "Failed to encode response",
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(response),
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	}, nil
}
