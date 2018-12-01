package main

import (
	"app/model"
	"app/shared/parse"
	"encoding/json"
	"log"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/headzoo/surf.v1"
)

var (
	_bcpa    model.Bcpa
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
	_bcpa = parse.LoadBcpaFromDoc(doc)

	//Load the class level BCPA object with with assessments
	parse.LoadAppendPropertyAssessments(doc, &_bcpa)

	//load exemptions
	parse.LoadAppendExemptionsTaxable(doc, &_bcpa)

	//Load Sales History
	parse.LoadSalesHistory(doc, &_bcpa)

	//Load the Land Calculations
	parse.LoadLandCalculations(doc, &_bcpa)

	//Load the Special Assessments
	parse.LoadSpecialAssessments(doc, &_bcpa)

	//Check if we have a URL for the CARD page. If so Parse it for data
	if len(_bcpa.LandCalculations.Cards) > 0 {

		//Let's loop the cards if more then one
		i := 0
		for _, card := range _bcpa.LandCalculations.Cards {

			//Grab the URL from the card
			CardURL, error := url.QueryUnescape(card.CardURL)

			if error != nil {
				log.Fatal(error)
			}

			//Start parseing the page
			error = parse.ExtractCardURL(CardURL, i, &_bcpa)

			//Increment
			i++

			if error != nil {
				log.Fatal(error)
			}

		}
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(parse.MarshalBcpa(_bcpa)),
		Headers: map[string]string{
			"Content-Type": "text/json",
		},
	}, nil

}

func main() {
	lambda.Start(Handler)
}
