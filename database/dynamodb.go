package database

import (
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

type DynamoDB struct {
	client    *dynamodb.DynamoDB
	tableName string
}

func NewDynamoDB() *DynamoDB {
	sess := session.Must(session.NewSession())
	return &DynamoDB{
		client:    dynamodb.New(sess),
		tableName: os.Getenv("TABLE_NAME"),
	}
}

func (d *DynamoDB) StoreURL(shortURL, originalURL string) error {
	item := map[string]*dynamodb.AttributeValue{
		"short_url":    {S: aws.String(shortURL)},
		"original_url": {S: aws.String(originalURL)},
	}
	_, err := d.client.PutItem(&dynamodb.PutItemInput{
		TableName: aws.String(d.tableName),
		Item:      item,
	})
	return err
}

func (d *DynamoDB) GetURL(shortURL string) (string, error) {
	result, err := d.client.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String(d.tableName),
		Key: map[string]*dynamodb.AttributeValue{
			"short_url": {S: aws.String(shortURL)},
		},
	})
	if err != nil || result.Item == nil {
		return "", err
	}
	return *result.Item["original_url"].S, nil
}
