package main

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/gorilla/mux"
	"net/http"
	"os"
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

//get
func GetDevices(w http.ResponseWriter, r *http.Request) {
	var retValue int
    vars := mux.Vars(r)
    retValue,device = getItemsDB(vars["id"])
    if retValue == 0{
    	w.WriteHeader(404)
    	fmt.Fprint(w, "Not Found")
    }else if retValue ==2{
    	w.WriteHeader(500)
    	fmt.Fprint(w, "Internal Server Error")
    }else{
	w.Header().Set("Content-Type", "/device/;   charset=UTF-8")
	w.WriteHeader(200) // unprocessable entity
	fmt.Fprintln(w, "OK")
	json.NewEncoder(w).Encode(device)
}
}

//post
func PostDevices(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Bad Request", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&device)
	if err != nil {
		http.Error(w, "Internal Server Error", 500)
		http.Error(w, err.Error(), 500)
		return
	}
	w.Header().Set("Content-Type", "/devices; charset=UTF-8")
	w.WriteHeader(201) // unprocessable entity
	fmt.Fprintln(w, "Created.")
	getConnectDB(device)
}

// Read items in database(0: request error 1: get items 2:server errors)
func getItemsDB(itemID string) (int,Device) {
	
	sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-1")},
    )
    // Create DynamoDB client
    svc := dynamodb.New(sess)

    result, err := svc.GetItem(&dynamodb.GetItemInput{
        TableName: aws.String("Test"),
        Key: map[string]*dynamodb.AttributeValue{
            "id": {
                N: aws.String(itemID),
            },
        },
    })

    if err != nil {
    	    fmt.Println(err.Error())
        return 2,device
    }

    device := Device{}

    err = dynamodbattribute.UnmarshalMap(result.Item, &device)

    if err != nil {
        panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
        return 2,device
    }
    if device.Id == 0 {
        fmt.Println("Could not find the item")
        return 0, device
    }
    return 1,device
}

// Create items in database
func getConnectDB(item Device){
	
	sess, err := session.NewSession(&aws.Config{
        Region: aws.String("us-east-1")},
    )
    // Create DynamoDB client
    svc := dynamodb.New(sess)

    sensor, err := dynamodbattribute.MarshalMap(item)

    if err != nil {
        fmt.Println("Got error marshalling map:")
        fmt.Println(err.Error())
        os.Exit(1)
    }

    // Create item in table Movies
    input := &dynamodb.PutItemInput{
        Item: sensor,
        TableName: aws.String("Test"),
    }

    _, err = svc.PutItem(input)

    if err != nil {
        fmt.Println("Got error calling PutItem:")
        fmt.Println(err.Error())
        os.Exit(1)
    }
    fmt.Println("Successfully added items to Device table!")
}


func main() {
	
	router := mux.NewRouter()
	router.HandleFunc("/device/{id}", GetDevices).Methods("GET")
	router.HandleFunc("/devices", PostDevices).Methods("POST")
	log.Fatal(http.ListenAndServe(":8080", router))
}
