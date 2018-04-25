package main

import (
	
	"context"
    "encoding/json"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// device type
type Device struct {
	Id          int    `json:"id"`
	DeviceModel string `json:"deviceModel"`
	Name        string `json:"name"`
	Note        string `json:"note"`
	Serial      string `json:"serial"`
}


var device Device

// to connect the database and to use get method to get the device with id from database
func getDevices(ctx context.Context, request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse, error){
	var itemID = request.PathParameters["id"]
	sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-1")},
    )
	
	if err != nil{
		return events.APIGatewayProxyResponse{
    	    	    Body: "Internal Server Error,"+err.Error(),
    	    	    StatusCode: 500,
    	    },nil
	}
	
    // Create DynamoDB client
    svc := dynamodb.New(sess)

// get the item with id
    result, err := svc.GetItem(&dynamodb.GetItemInput{
        TableName: aws.String("Test"),
        Key: map[string]*dynamodb.AttributeValue{
            "id": {
                N: aws.String(itemID),
            },
        },
    })

    if err != nil {
    	    return events.APIGatewayProxyResponse{
    	    	    Body: err.Error(),
    	    	    StatusCode: 404,
    	    },nil
    }

// unmarshal the data from database
    err = dynamodbattribute.UnmarshalMap(result.Item, &device)
    
    if err != nil {
        return events.APIGatewayProxyResponse{
        	    Body: "Internal server error,"+err.Error(),
    	    	    StatusCode: 500,
        },nil
    }else if device.Id == 0 {
    	    return events.APIGatewayProxyResponse{
    	    	    Body: "Could not find the item!",
    	    	    StatusCode: 404,
    	    },nil
    }else{
    	    outPut, err := json.Marshal(device)
        if err != nil {
            panic (err)
         }
        return events.APIGatewayProxyResponse{
        	    Body: "OK,"+string(outPut),
        	    StatusCode: 200,
        },nil
    }   
}

func main() {
	lambda.Start(getDevices)
}