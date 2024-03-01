package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"errors"
)

func main() {
	lambda.Start(handler)
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if request.HTTPMethod == "GET" {
		return getFriendsHandler(request)
	} else if request.HTTPMethod == "POST" {
		return addFriendHandler(request)
	} 

	return events.APIGatewayProxyResponse{StatusCode: 500}, errors.New("Invalid HTTP method: " + request.HTTPMethod)
	// return events.APIGatewayProxyResponse{Body: "woo", StatusCode: 200}, nil
}

