package main

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
	guuid "github.com/google/uuid"
	"time"
)

type LambdaRequest struct {
	Data string `json:"data"`
}

func CreateMessage(request events.APIGatewayWebsocketProxyRequest) (events.APIGatewayProxyResponse, error) {
	lambdaRequest := LambdaRequest{}
	json.Unmarshal([]byte(request.Body), &lambdaRequest)

	message := buildMessage(lambdaRequest.Data)

	_, err := SaveMessage(dbSession(), message)
	if err != nil {
		returnFormattedError(err)
	}

	body, err := json.Marshal(message)
	if err != nil {
		returnFormattedError(err)
	}

	return events.APIGatewayProxyResponse{Body: string(body), StatusCode: 201}, nil
}

func dbSession() (dbSession dynamodbiface.DynamoDBAPI) {
	sess := session.Must(session.NewSession())

	return dynamodb.New(sess)
}

func returnFormattedError(err error) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{Body: err.Error(), StatusCode: 422}, nil
}

func buildMessage(data string) (message Message) {
	return Message{
		PrimaryKey:		fmt.Sprintf("MESSAGE_%s", guuid.New().String()),
		SecondaryKey:	fmt.Sprintf("CHAT_%s", "42"),
		CreatedAt: 		time.Now().Format(time.RFC3339),
		Data:      		data,
	}
}
