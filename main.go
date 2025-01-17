package main

import (
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"url-shortener/database"
	"url-shortener/handlers"
)

var dbClient database.Database
var isLambda = os.Getenv("AWS_EXECUTION_ENV") != ""

func main() {
	if isLambda {
		// Use DynamoDB for Lambda
		dbClient = database.NewDynamoDB()

		// Start Lambda with API Gateway Handler
		lambda.Start(func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
			log.Printf("Request type: %v", request.HTTPMethod)

			switch request.HTTPMethod {
			case "POST":
				return handlers.StoreURLHandlerAPI(dbClient, request.Body)
			case "GET":
				shortURL := request.PathParameters["short_url"]
				return handlers.RedirectURLHandlerAPI(dbClient, shortURL)
			default:
				return events.APIGatewayProxyResponse{
					StatusCode: 405,
					Body:       "Method Not Allowed",
				}, nil
			}
		})
	} else {
		// Use SQLite for local development
		dbClient = database.NewSQLiteDB()

		// Setup HTTP routes for local testing
		http.HandleFunc("/shorten", handlers.StoreURLHandler(dbClient))
		http.HandleFunc("/", handlers.RedirectURLHandler(dbClient))

		log.Println("Server running on :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}
}
