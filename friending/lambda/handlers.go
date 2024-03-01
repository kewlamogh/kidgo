package main

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"

	// "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/kewlamogh/kidgo-backend/friending"

	"log"
)

func getFriendsHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("recieved request")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Minute)
	config, err := config.LoadDefaultConfig(ctx)
	defer cancelFunc()

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	dynamo := dynamodb.NewFromConfig(config)
	client := friending.DynamoClient{
		DynamoDBClient: dynamo,
		TableName:      "friends",
	}

	stuff, err := client.GetFriendList(request.QueryStringParameters["username"])
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       strings.Join(stuff, ","),
	}

	return response, nil
}

func addFriendHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Println("recieved request")

	ctx, cancelFunc := context.WithTimeout(context.Background(), 5*time.Minute)
	config, err := config.LoadDefaultConfig(ctx)
	defer cancelFunc()

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	dynamo := dynamodb.NewFromConfig(config)
	client := friending.DynamoClient{
		DynamoDBClient: dynamo,
		TableName:      "friends",
	}

	err = client.AddFriend(request.QueryStringParameters["username"], request.QueryStringParameters["friend"])
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	response := events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       "the job is done",
	}

	return response, nil
}
