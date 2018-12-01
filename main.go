package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"gopkg.in/headzoo/surf.v1"
)

var (
	_bcpa    Bcpa
	_baseURL string = "http://www.bcpa.net/"
)

// GenericError base error message
type GenericError struct {
	Message string `json:"message"`
	Code    string `json:"code"`
	Object  string `json:"object"`
}

// Bcpa table contains the information for each user
type Bcpa struct {
	Siteaddress         string `json:"siteaddress"`
	Owner               string `json:"owner"`
	MailingAddress      string `json:"mailingAddress"`
	ID                  string `json:"id"`
	Milage              string `json:"milage"`
	Use                 string `json:"use"`
	Legal               string `json:"legal"`
	PropertyAssessments []PropertyAssessmentValue
	ExemptionsTaxable   ExemptionsTaxableValuesbyTaxingAuthority
	SalesHistory        []Sale
	LandCalculations    LandCalculations
	SpecialAssessments  []SpecialAssessment
}

// RecBuildingCard
type RecBuildingCard struct {
	CardURL                   string `json:"cardurl"`
	TaxYear                   string `json:"taxyear"`
	Folio                     string `json:"folio"`
	ParcelIDNumber            string `json:"parcelidnumber"`
	UseCode                   string `json:"usecode"`
	NoBedrooms                string `json:"nobedrooms"`
	NoBaths                   string `json:"nobaths"`
	NoUnits                   string `json:"nounits"`
	NoStories                 string `json:"nostories"`
	NoBuildings               string `json:"nobuildings"`
	Foundation                string `json:"foundation"`
	Exterior                  string `json:"exterior"`
	RoofType                  string `json:"rooftype"`
	RoofMaterial              string `json:"roofmaterial"`
	Interior                  string `json:"interior"`
	Floors                    string `json:"floors"`
	Plumbing                  string `json:"plumbing"`
	Electric                  string `json:"electric"`
	Classification            string `json:"classification"`
	CeilingHeights            string `json:"ceilingheights"`
	QualityOfConstruction     string `json:"qualityofconstruction"`
	CurrentConditionStructure string `json:"currentconditionstructure"`
	ConstructionClass         string `json:"constructionclass"`
	Permits                   []Permit
	ExtraFeatures             []ExtraFeature
}

// ExtraFeature
type ExtraFeature struct {
	Feature string `json:"feature"`
}

// Permit
type Permit struct {
	PermitNo   string `json:"permitco"`
	PermitType string `json:"permittype"`
	EstCost    string `json:"estcost"`
	PermitDate string `json:"permitdate"`
	CODate     string `json:"codate"`
}

// LandCalculations
type LandCalculations struct {
	Calculations    []LandCalculation
	AdjBldgSF       string `json:"adjbldgsf"`
	Units           string `json:"units"`
	Cards           []RecBuildingCard
	SketchURL       string `json:"sketchurl"`
	EffActYearBuilt string `json:"effactyearbuilt"`
}

// LandCalculation
type LandCalculation struct {
	Price  string `json:"price"`
	Factor string `json:"factor"`
	Type   string `json:"type"`
}

// SpecialAssessment data
type SpecialAssessment struct {
	Fire  string `json:"fire"`
	Garb  string `json:"garb"`
	Light string `json:"light"`
	Drain string `json:"drain"`
	Impr  string `json:"impr"`
	Safe  string `json:"safe"`
	Storm string `json:"storm"`
	Clean string `json:"clean"`
	Misc  string `json:"misc"`
}

// RecPatriotSketch
type RecPatriotSketch struct {
	Sketch       string `json:"sketch"`
	Building     string `json:"building"`
	URL          string `json:"url"`
	SketchImgURL string `json:"sketchimgurl"`
	Codes        []PatriotSketchCode
	AdjAreaTotal string `json:"adjareatotal"`
}

// PatriotSketchCode
type PatriotSketchCode struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Area        string `json:"area"`
	Factor      string `json:"factor"`
	AdjArea     string `json:"adjarea"`
	Stories     string `json:"stories"`
}

// ExemptionsTaxableValuesbyTaxingAuthority table contains the exemptions
type ExemptionsTaxableValuesbyTaxingAuthority struct {
	County      ExemptionsAndTaxableValue
	SchoolBoard ExemptionsAndTaxableValue
	Municipal   ExemptionsAndTaxableValue
	Independent ExemptionsAndTaxableValue
	CreatedAt   time.Time `json:"createdat"`
	UpdatedAt   time.Time `json:"updatedat"`
}

// ExemptionsAndTaxableValue table contains the exemption values
type ExemptionsAndTaxableValue struct {
	JustValue    string `json:"justvalue"`
	Portability  string `json:"portability"`
	AssessedSOH  string `json:"assessedsoh"`
	Homestead    string `json:"homestead"`
	AddHomestead string `json:"addhomestead"`
	WidVetDis    string `json:"widvetdis"`
	Senior       string `json:"senior"`
	XemptType    string `json:"xempttype"`
	Taxable      string `json:"taxable"`
}

// PropertyAssessmentValue table contains the house values
type PropertyAssessmentValue struct {
	Year                string    `json:"year"`
	Land                string    `json:"land"`
	BuildingImprovement string    `json:"buildingimprovement"`
	JustMarketValue     string    `json:"justmarketvalue"`
	AssessedSOHValue    string    `json:"assessedsohvalue"`
	Tax                 string    `json:"tax"`
	CreatedAt           time.Time `json:"createdat"`
	UpdatedAt           time.Time `json:"updatedat"`
}

// Sale property sales
type Sale struct {
	Date        string `json:"date"`
	Type        string `json:"type"`
	Price       string `json:"price"`
	BookPageCIN string `json:"bookpagecin"`
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
func LoadBcpaFromDoc(doc *goquery.Document) Bcpa {

	var bcpa Bcpa
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
func marshalBcpa(bcpa Bcpa) string {
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
