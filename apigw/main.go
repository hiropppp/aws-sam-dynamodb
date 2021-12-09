package main

import (
	// "context"
	"encoding/json"
	"fmt"
	// "net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var db *dynamodb.DynamoDB
var dbEndpoint = "http://test_dynamodb-local:8000"
var region = "ap-northeast-1"
var testTable = "Master_Mangement"

type Iyaku struct {
	MangementID string `json:"mangementid"`
	Upload_File_Name    string `json:"upload_file_name"`
	Status string `json:"status"`
	Upload_Num string `json:"upload_num"`
	Upload_Date string `json:"upload_date"`
	FileVersion string `json:"fileversion"`
	Search_Value string `json:"search_value"`
	Upload_Master_Name string `json:"upload_master_name"`
	Upload_Person_Name string `json:"upload_person_name"`
	Csv_URL string `json:"csv_url"`
}

func handler() (events.APIGatewayProxyResponse, error) {

	sess := session.Must(session.NewSession(&aws.Config{
		Endpoint: aws.String(dbEndpoint),
		Region:   aws.String(region),
	}))
	db = dynamodb.New(sess)

	// input := &dynamodb.QueryInput{
    //     TableName: aws.String(testTable),
    //     ExpressionAttributeNames: map[string]*string{
    //         "#COM":   aws.String("MangementID"), 
    //     },
    //     ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
    //         ":com": { // :を付けるのがセオリーのようです
    //             S: aws.String("mas-iyaku"),
    //         },
    //     },
    //     KeyConditionExpression: aws.String("#COM = :com"),         // 検索条件
    //     ProjectionExpression:   aws.String("#COM"), // 取得カラム
    //     ScanIndexForward:       aws.Bool(true),                 // ソートキーのソート順（指定しないと昇順）
    //     Limit:                  aws.Int64(20),                  // 取得件数の指定もできる
    // }

	input1 := &dynamodb.QueryInput{
        TableName: aws.String(testTable),
		IndexName: aws.String("Search_Value-Upload_Date-index"),
        // ExpressionAttributeNames: map[string]*string{
        //     "#Search":   aws.String("Search_Value"), 
		// 	"#Date":     aws.String("Upload_Date"), 
        // },
        ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
            ":search": { // :を付けるのがセオリーのようです
                S: aws.String("dummy"),
            },
        },
        KeyConditionExpression: aws.String("Search_Value = :search"),         // 検索条件
        ProjectionExpression:   aws.String("MangementID, Search_Value, Upload_Date"), // 取得カラム
        ScanIndexForward:       aws.Bool(false),                 // ソートキーのソート順（指定しないと昇順）
        Limit:                  aws.Int64(30),                  // 取得件数の指定もできる
    }

	input2 := &dynamodb.QueryInput{
        TableName: aws.String(testTable),
		IndexName: aws.String("Search_Value-Upload_Date-index"),
        // ExpressionAttributeNames: map[string]*string{
        //     "#Search":   aws.String("Search_Value"), 
		// 	"#Date":     aws.String("Upload_Date"), 
        // },
        ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
            ":search": { // :を付けるのがセオリーのようです
                S: aws.String("master"),
            },
        },
        KeyConditionExpression: aws.String("Search_Value = :search"),         // 検索条件
        ProjectionExpression:   aws.String("MangementID, Search_Value, Upload_Date"), // 取得カラム
        ScanIndexForward:       aws.Bool(false),                 // ソートキーのソート順（指定しないと昇順）
        Limit:                  aws.Int64(30),                  // 取得件数の指定もできる
    }

	result1, err := db.Query(input1)
    if err != nil {
        fmt.Println("[Query Error]", err)
		// return events.APIGatewayProxyResponse{
		// 	Body:       err,
		// 	StatusCode: 500,
		// }, nil
    }

	result2, err := db.Query(input2)
    if err != nil {
        fmt.Println("[Query Error]", err)
		// return events.APIGatewayProxyResponse{
		// 	Body:       err,
		// 	StatusCode: 500,
		// }, nil
    }

	history1 := make([]*Iyaku, 0)
    if err := dynamodbattribute.UnmarshalListOfMaps(result1.Items, &history1); err != nil {
        fmt.Println("[Unmarshal Error]", err)
		// return events.APIGatewayProxyResponse{
		// 	Body:       err,
		// 	StatusCode: 500,
		// }, nil
    }

	history2 := make([]*Iyaku, 0)
    if err := dynamodbattribute.UnmarshalListOfMaps(result2.Items, &history2); err != nil {
        fmt.Println("[Unmarshal Error]", err)
		// return events.APIGatewayProxyResponse{
		// 	Body:       err,
		// 	StatusCode: 500,
		// }, nil
    }


	fmt.Println(result1)

	j, _ := json.Marshal(history1)
	k, _ := json.Marshal(history2)

	// var m map[string]interface{}
	// json.Unmarshal(j,&m)
	// json.Unmarshal(k,&m)

	// fmt.Println(m)
	// i, _ := json.Marshal(m)


    // fmt.Println(string(j))
	// fmt.Println(string(k))

	// w.Write(j)
	return events.APIGatewayProxyResponse{
		Body:       string(k) + string(j),
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}