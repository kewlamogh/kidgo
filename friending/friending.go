package friending

import (
	"context"
	// "log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"golang.org/x/exp/slices"
)

type DynamoClient struct {
	DynamoDBClient *dynamodb.Client
	TableName string
}

func (client DynamoClient) GetFriendList(username string) ([]string, error) {
	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute)
	defer cancelFunc()

	marshalled, err := attributevalue.Marshal(username)

	if err != nil {
		return nil, err
	}

	user, err := client.DynamoDBClient.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(client.TableName),
		Key:       map[string]types.AttributeValue{"username": marshalled},
	})

	var ss []string
	err = attributevalue.Unmarshal(user.Item["friends"], &ss)
	if err != nil {
		return nil, err
	}

	return ss, nil
}

func (client DynamoClient) AddFriend(username string, friend string) error {
	existing_friends, err := client.GetFriendList(username)
	if err != nil {
		return err
	}

	ctx, cancelFunc := context.WithTimeout(context.TODO(), time.Minute)
	defer cancelFunc()

	uname, err := attributevalue.Marshal(username)
	if err != nil {
		return err
	}

	if !slices.Contains(existing_friends, friend) {
		existing_friends = append(existing_friends, friend)
	} else {
		return nil
	}

	friends, err := attributevalue.Marshal(existing_friends)

	if err != nil {
		return err
	}

	update := expression.Set(expression.Name("friends"), expression.Value(friends))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return err
	}

	_, err = client.DynamoDBClient.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(client.TableName),
		Key: map[string]types.AttributeValue{"username": uname},
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression: expr.Update(),
	})
	if err != nil {
		return err
	}

	return nil
}