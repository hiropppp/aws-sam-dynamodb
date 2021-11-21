package main

import (
	"errors"
	"fmt"
	"time"
	"io/ioutil"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

var (
	// DefaultHTTPGetAddress Default Address
	DefaultHTTPGetAddress = "https://checkip.amazonaws.com"

	// ErrNoIP No IP found in response
	ErrNoIP = errors.New("No IP in HTTP response")

	// ErrNon200Response non 200 status code in response
	ErrNon200Response = errors.New("Non 200 Response found")
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	resp, err := http.Get(DefaultHTTPGetAddress)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if resp.StatusCode != 200 {
		return events.APIGatewayProxyResponse{}, ErrNon200Response
	}

	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	if len(ip) == 0 {
		return events.APIGatewayProxyResponse{}, ErrNoIP
	}

	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf("Hello, %v", string(ip)),
		StatusCode: 200,
	}, nil
}

func main() {

	// 1ヶ月後
	t1 := time.Date(2020, 1, 31, 23, 59, 59, 999999999, time.UTC)
	fmt.Println(t1.AddDate(0,1,0))

	// 1年後
	t2 := time.Date(2019, 2, 28, 23, 59, 59, 999999999, time.UTC)
	fmt.Println(t2.AddDate(1,0,0))

	// 1年後
	t3 := time.Date(2020, 2, 29, 23, 59, 59, 999999999, time.UTC)
	fmt.Println(t3.AddDate(1,0,0))

	lambda.Start(handler)

}
