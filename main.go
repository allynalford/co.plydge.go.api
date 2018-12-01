package main

import (
	"app/model"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/headzoo/surf.v1"
)

var (
	_bcpa   model.Bcpa{}
	_baseURL = "http://www.bcpa.net/"
)

// GenericError base error message
type GenericError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Object  string `json:"object"`
}

// ************************************************************************ Functions

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

// GenericAPIProxyResponse function to create base response message with events.APIGatewayProxyResponse
func GenericAPIProxyResponse(c int, b string, h map[string]string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: c,
		Body:       b,
		Headers:    h,
	}, nil
}

//LoadBcpaFromDoc used to load Bcpa data from HTML
func LoadBcpaFromDoc(doc *goquery.Document) model.Bcpa {

	bcpa := model.Bcpa{}
	var siteAddress, owner, mailingAddress, id, mileage, use, legal string

	// use selector found with the browser inspector
	siteAddress = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(1) > table > tbody > tr:nth-child(1) > td:nth-child(2) > span > a > b").Contents().Text()

	//clean up the carriage return
	re := regexp.MustCompile(`\r?\n`)
	siteAddress = re.ReplaceAllString(siteAddress, " ")

	//Set the Object
	bcpa.Siteaddress = strings.TrimSpace(StripSpaces(siteAddress))

	owner = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(1) > table > tbody > tr:nth-child(2) > td:nth-child(2) > span").Contents().Text()
	//Set the Object
	bcpa.Owner = strings.TrimSpace(owner)

	mailingAddress = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(1) > table > tbody > tr:nth-child(3) > td:nth-child(2) > span").Contents().Text()

	//Set the Object
	bcpa.MailingAddress = strings.TrimSpace(mailingAddress)

	id = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(3) > table > tbody > tr:nth-child(1) > td:nth-child(2) > span").Contents().Text()

	//Set the Object
	bcpa.ID = strings.TrimSpace(strings.Replace(StripSpaces(id), " ", "", -1))

	mileage = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(3) > table > tbody > tr:nth-child(2) > td:nth-child(2) > span").Contents().Text()

	//Set the Object
	bcpa.Milage = strings.TrimSpace(mileage)

	use = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(2) > tbody > tr > td:nth-child(3) > table > tbody > tr:nth-child(3) > td:nth-child(2) > span").Contents().Text()

	//Set the Object
	bcpa.Use = strings.TrimSpace(StripSpaces(use))

	legal = doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(4) > tbody > tr > td:nth-child(2) > span").Contents().Text()

	//Set the Object
	bcpa.Legal = strings.TrimSpace(legal)

	return bcpa
}

//StripSpaces remove leading and trailing and extra gapped spaces
func StripSpaces(o string) string {

	releadclosewhtsp2 := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)

	reinsidewhtsp2 := regexp.MustCompile(`[\s\p{Zs}]{2,}`)

	final := releadclosewhtsp2.ReplaceAllString(o, "")

	return reinsidewhtsp2.ReplaceAllString(final, " ")
}

// marshalBcpa Convert BCPA	to string
func marshalBcpa(bcpa model.Bcpa) string {
	//user := &User{name:"Frank"}
	b, err := json.Marshal(bcpa)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return "0"
	}
	fmt.Println(string(b))

	return string(b)
}

// Handler is executed by AWS Lambda in the main function. Once the request
// is processed, it returns an Amazon API Gateway response object to AWS Lambda
func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// Create a new browser and open reddit.
	bow := surf.NewBrowser()
	err := bow.Open(_baseURL + "RecAddr.asp")

	//Ensure no error opening page
	if err != nil {
		return GenerateErrorResponse(err.Error()+" - Error while opening: "+_baseURL+"RecAddr.asp ", "1", "")
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

	//Load the BCPA parent node from the HTML receieved from URL
	_bcpa = LoadBcpaFromDoc(doc)

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(marshalBcpa(_bcpa)),
		Headers: map[string]string{
			"Content-Type": "text/json",
		},
	}, nil

}

func main() {
	lambda.Start(Handler)
}
