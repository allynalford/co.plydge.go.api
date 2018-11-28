package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/headzoo/surf.v1"
)

// GenericError base error message
type GenericError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Object  string `json:"object"`
}

// GenerateError function to create base error message
func GenerateError(m string, c string, o string) GenericError {
	genericError := GenericError{m, c, o}
	return genericError
}

// GenerateErrorString function to create base error message
func GenerateErrorString(m string, c string, o string) ([]byte, error) {
	genericError := GenerateError(m, c, o)
	ge, err := json.Marshal(genericError)
	return ge, err
}

// GenerateErrorResponse function to create base error message with events.APIGatewayProxyResponse
func GenerateErrorResponse(m string, c string, o string) (events.APIGatewayProxyResponse, error) {
	ge, err := GenerateErrorString(m, c, o)

	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(ge),
		Headers: map[string]string{
			"Content-Type": "text/json",
		},
	}, nil
}

// GenericApiProxyResponse function to create base response message with events.APIGatewayProxyResponse
func GenericApiProxyResponse(c int, b string, h map[string]string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: c,
		Body:       b,
		Headers:    h,
	}, nil
}

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Create a new browser and open reddit.
	bow := surf.NewBrowser()
	err := bow.Open("http://www.bcpa.net/RecAddr.asp")

	//Ensure no error opening page
	if err != nil {
		return GenerateErrorResponse(err.Error()+" - Error while opening: http://www.bcpa.net/RecAddr.asp ", "1", "")
	}

	if len(request.QueryStringParameters) < 6 {

		return GenerateErrorResponse("Parameters: Invalid Parameter Length", "2", "")
	}

	SitusStreetNumber, ok := request.QueryStringParameters["SN"]

	if !ok {
		return GenerateErrorResponse("Parameters:Missing Street Number", "3", SitusStreetNumber)
	}

	SitusUnitNumber, ok := request.QueryStringParameters["UN"]

	if !ok {
		return GenerateErrorResponse("Parameters: Missing Unit Number", "4", SitusUnitNumber)
	}

	SitusStreetDirection, ok := request.QueryStringParameters["SD"]

	if !ok {
		return GenerateErrorResponse("Parameters: Missing Street Direction", "5", SitusStreetDirection)
	}

	SitusStreetName, ok := request.QueryStringParameters["HN"]

	if !ok {
		return GenerateErrorResponse("Parameters: Missing Street Name", "6", SitusStreetName)
	}

	SitusStreetType, ok := request.QueryStringParameters["ST"]

	if !ok {
		return GenerateErrorResponse("Parameters: Missing Street Type", "7", SitusStreetType)
	}

	SitusStreetPostDir, ok := request.QueryStringParameters["PD"]

	if !ok {
		return GenerateErrorResponse("Parameters: Missing Street Post Direction", "8", SitusStreetPostDir)
	}

	City, ok := request.QueryStringParameters["CT"]

	if !ok {
		return GenerateErrorResponse("Parameters: Missing City", "9", City)
	}

	// Submit the search form
	fm, _ := bow.Form("[name='homeind']")
	fm.Input("Situs_Street_Number", SitusStreetNumber)
	fm.SelectByOptionValue("Situs_Street_Direction", SitusStreetDirection)
	fm.Input("Situs_Street_Name", SitusStreetName)
	fm.SelectByOptionValue("Situs_Street_Type", SitusStreetType)
	fm.Input("Situs_Street_Post_Dir", SitusStreetPostDir)
	fm.Input("Situs_Unit_Number", SitusUnitNumber)
	fm.SelectByOptionValue("Situs_City", City)

	if fm.Submit() != nil {
		return GenerateErrorResponse(err.Error(), "1.1", "")
	}

	doc, err := goquery.NewDocument(bow.Url().String())
	if err != nil {
		return GenerateErrorResponse(err.Error(), "1.1", "")
	}

	var result, siteAddress, owner, mailingAddress, id, mileage, use, legal string

	// use selector found with the browser inspector
	siteAddress = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(1) > table > tbody > tr:nth-child(1) > td:nth-child(2) > span > a > b").Contents().Text()

	//clean up the carriage return
	re := regexp.MustCompile(`\r?\n`)
	siteAddress = re.ReplaceAllString(siteAddress, " ")
	//siteAddress = strings.Replace(siteAddress, " 			  ", " ", 1)

	releadclosewhtsp := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	reinsidewhtsp := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	final := releadclosewhtsp.ReplaceAllString(siteAddress, "")
	siteAddress = reinsidewhtsp.ReplaceAllString(final, " ")

	//trim whitespace
	siteAddress = strings.TrimSpace(siteAddress)

	owner = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(1) > table > tbody > tr:nth-child(2) > td:nth-child(2) > span").Contents().Text()
	owner = strings.TrimSpace(owner)

	mailingAddress = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(1) > table > tbody > tr:nth-child(3) > td:nth-child(2) > span").Contents().Text()
	mailingAddress = strings.TrimSpace(mailingAddress)

	id = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(3) > table > tbody > tr:nth-child(1) > td:nth-child(2) > span").Contents().Text()

	releadclosewhtsp2 := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)
	reinsidewhtsp2 := regexp.MustCompile(`[\s\p{Zs}]{2,}`)
	final2 := releadclosewhtsp2.ReplaceAllString(id, "")
	id = reinsidewhtsp2.ReplaceAllString(final2, " ")

	//trim
	id = strings.TrimSpace(id)

	mileage = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(3) > table > tbody > tr:nth-child(2) > td:nth-child(2) > span").Contents().Text()
	mileage = strings.TrimSpace(mileage)

	use = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(3) > table > tbody > tr:nth-child(3) > td:nth-child(2) > span").Contents().Text()
	use = strings.TrimSpace(mileage)

	legal = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(4) > tbody > tr > td:nth-child(2) > span").Contents().Text()
	legal = strings.TrimSpace(legal)

	result = fmt.Sprintf("{\"siteaddress\": \"%s\", \"owner\": \"%s\", \"mailingAddress\": \"%s\", \"id\": \"%s\", \"milage\": \"%s\", \"use\": \"%s\", \"legal\": \"%s\"}", siteAddress, owner, mailingAddress, id, mileage, use, legal)

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
