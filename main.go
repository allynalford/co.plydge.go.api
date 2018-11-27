package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/headzoo/surf.v1"

	"github.com/PuerkitoBio/goquery"
)

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//index, err := ioutil.ReadFile("public/index.html")
	//if err != nil {
	//	return events.APIGatewayProxyResponse{}, err
	//	}

	// Create a new browser and open reddit.
	bow := surf.NewBrowser()
	err := bow.Open("http://www.bcpa.net/RecAddr.asp")
	if err != nil {
		panic(err)
	}

	SitusStreetNumber := request.QueryStringParameters["SN"]
	SitusUnitNumber := request.QueryStringParameters["UN"]

	// Submit the search form
	fm, _ := bow.Form("[name='homeind']")
	fm.Input("Situs_Street_Number", SitusStreetNumber)
	fm.SelectByOptionValue("Situs_Street_Direction", "SW")
	fm.Input("Situs_Street_Name", "18")
	fm.SelectByOptionValue("Situs_Street_Type", "AVE")
	fm.Input("Situs_Street_Post_Dir", "")
	fm.Input("Situs_Unit_Number", SitusUnitNumber)
	fm.SelectByOptionValue("Situs_City", "FL")

	if fm.Submit() != nil {
		panic(err)
	}

	doc, err := goquery.NewDocument(bow.Url().String())
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	//fmt.Printf(doc.Html())
	var result string
	var siteAddress string
	var owner string

	// use CSS selector found with the browser inspector
	// for each, use index and item
	siteAddress = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(1) > table > tbody > tr:nth-child(1) > td:nth-child(2) > span > a > b").Contents().Text()
	siteAddress = strings.TrimSpace(siteAddress)

	owner = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(1) > table > tbody > tr:nth-child(2) > td:nth-child(2) > span").Contents().Text()
	owner = strings.TrimSpace(owner)

	result = fmt.Sprintf("{\"siteaddress\": \"%s\", \"owner\": \"%s\"}", siteAddress, owner)

	//fmt.Printf(result)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(result),
		Headers: map[string]string{
			"Content-Type": "text/json",
		},
	}, nil

}

func main() {
	lambda.Start(Handler)
}
