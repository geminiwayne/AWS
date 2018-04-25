package main

import (
	"strconv"
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

// Create items in database by http post
func PostDevices(ctx context.Context, request events.APIGatewayProxyRequest)(events.APIGatewayProxyResponse, error){
	sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-1")},
    )
	if err != nil{
		return events.APIGatewayProxyResponse{
    	    	    Body: "Internal Server Error,"+err.Error(),
    	    	    StatusCode: 500,
    	    },nil
	}
	
	device := &Device{}
    // Create DynamoDB client
    svc := dynamodb.New(sess)
    json.Unmarshal([]byte(request.Body),device)
    sensor,err := dynamodbattribute.MarshalMap(device)
    
    //check id validation
    result, getErr := svc.GetItem(&dynamodb.GetItemInput{
        TableName: aws.String("Test"),
        Key: map[string]*dynamodb.AttributeValue{
            "id": {
                N: aws.String(strconv.Itoa(device.Id)),
            },
        },
    })
    
    // check existed id
     if result.Item != nil{
       return events.APIGatewayProxyResponse{ // Error HTTP response
			Body:"Duplicated ID!",
			StatusCode: 400,
		}, nil
    }
     // invalid id
     if device.Id ==0{
       return events.APIGatewayProxyResponse{ // Error HTTP response
			Body:"Invalid input",
			StatusCode: 400,
		}, nil
     }
    // error in validation
    if getErr !=nil{
        return events.APIGatewayProxyResponse{ // Error HTTP response
			Body:"Bad Request!",
			StatusCode: 400,
		}, nil
    }
    // Create item in table device
    input := &dynamodb.PutItemInput{
        Item: sensor,
        TableName: aws.String("Test"),
    }
    if err != nil {
        return events.APIGatewayProxyResponse{
    	    	    Body: "Internal Server Error,"+err.Error(),
    	    	    StatusCode: 500,
    	    },nil
}
    if _, err := svc.PutItem(input); err !=nil{
    	    return events.APIGatewayProxyResponse{ // Error HTTP response
			Body: "Bad Request,"+err.Error(),
			StatusCode: 400,
		}, nil
    }else{
        return events.APIGatewayProxyResponse{ // success response
            Body: "Created",
			StatusCode: 201,
    },nil
}
}

func main() {
	lambda.Start(PostDevices)
}