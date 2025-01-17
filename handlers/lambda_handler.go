package handlers

import (
	"url-shortener/database"

	"github.com/aws/aws-lambda-go/events"
)

func LambdaHandler(db database.Database) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
		switch request.HTTPMethod {
		case "POST":
			return StoreURLHandlerAPI(db, request.Body)
		case "GET":
			shortURL := request.PathParameters["short_url"]
			return RedirectURLHandlerAPI(db, shortURL)
		default:
			return events.APIGatewayProxyResponse{StatusCode: 405, Body: "Method not allowed"}, nil
		}
	}
}
