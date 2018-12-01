package parse

import (
	"app/model"
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

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

// MarshalBcpa Convert BCPA	to string
func MarshalBcpa(bcpa model.Bcpa) string {
	//user := &User{name:"Frank"}
	b, err := json.Marshal(bcpa)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return "0"
	}
	fmt.Println(string(b))

	return string(b)
}

//StripSpaces remove leading and trailing and extra gapped spaces
func StripSpaces(o string) string {

	releadclosewhtsp2 := regexp.MustCompile(`^[\s\p{Zs}]+|[\s\p{Zs}]+$`)

	reinsidewhtsp2 := regexp.MustCompile(`[\s\p{Zs}]{2,}`)

	final := releadclosewhtsp2.ReplaceAllString(o, "")

	return reinsidewhtsp2.ReplaceAllString(final, " ")
}

// PropertyAssessmentRecord PropertyAssessments  table contains the information for each user Called by LoadAppendPropertyAssessments
func PropertyAssessmentRecord(s *goquery.Selection) model.PropertyAssessmentValue {
	p := model.PropertyAssessmentValue{}

	// Loop through each cell
	s.Find("td").Each(func(int int, s *goquery.Selection) {

		switch int {
		case 0:
			p.Year = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 1:
			p.Land = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 2:
			p.BuildingImprovement = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 3:
			p.JustMarketValue = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 4:
			p.AssessedSOHValue = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 5:
			p.Tax = strings.TrimSpace(s.Find("span").First().Contents().Text())
		}
	})

	return p
}

//LoadAppendPropertyAssessments used to load and append Assessments to the BCPA parent node calls PropertyAssessmentRecord
func LoadAppendPropertyAssessments(doc *goquery.Document, _bcpa *model.Bcpa) {

	doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(6) > tbody > tr").Each(func(i int, s *goquery.Selection) {

		if i > 1 {
			pa := PropertyAssessmentRecord(s)
			pa.CreatedAt = time.Now()
			*_bcpa.PropertyAssessments = append(&_bcpa.PropertyAssessments, pa)
		}
	})

}

// ExemptionsTaxableRecord  parse exemptions called by LoadAppendExemptionsTaxable
func ExemptionsTaxableRecord(s *goquery.Selection, i int, eta model.ExemptionsTaxableValuesbyTaxingAuthority) model.ExemptionsTaxableValuesbyTaxingAuthority {

	// Loop through each cell
	s.Find("td").Each(func(int int, s *goquery.Selection) {

		switch i {
		case 2:
			switch int {
			case 1:
				eta.County.JustValue = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.JustValue = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.JustValue = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.JustValue = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		case 3:
			switch int {
			case 1:
				eta.County.Portability = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.Portability = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.Portability = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.Portability = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		case 4:
			switch int {
			case 1:
				eta.County.AssessedSOH = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.AssessedSOH = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.AssessedSOH = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.AssessedSOH = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		case 5:
			switch int {
			case 1:
				eta.County.Homestead = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.Homestead = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.Homestead = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.Homestead = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		case 6:
			switch int {
			case 1:
				eta.County.AddHomestead = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.AddHomestead = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.AddHomestead = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.AddHomestead = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		case 7:
			switch int {
			case 1:
				eta.County.WidVetDis = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.WidVetDis = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.WidVetDis = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.WidVetDis = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		case 8:
			switch int {
			case 1:
				eta.County.Senior = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.Senior = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.Senior = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.Senior = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		case 9:
			switch int {
			case 1:
				eta.County.XemptType = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.XemptType = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.XemptType = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.XemptType = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		case 10:
			switch int {
			case 1:
				eta.County.Taxable = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 2:
				eta.SchoolBoard.Taxable = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 3:
				eta.Municipal.Taxable = strings.TrimSpace(s.Find("span").First().Contents().Text())
			case 4:
				eta.Independent.Taxable = strings.TrimSpace(s.Find("span").First().Contents().Text())
			}
		}

	})

	return eta
}

// LoadAppendExemptionsTaxable Load Taxable and Exemptions Calls ExemptionsTaxableRecord
func LoadAppendExemptionsTaxable(doc *goquery.Document, _bcpa *model.Bcpa) {

	//Preload the object
	eta := model.ExemptionsTaxableValuesbyTaxingAuthority{}
	eta.CreatedAt = time.Now()
	eta.County = model.ExemptionsAndTaxableValue{}
	eta.SchoolBoard = model.ExemptionsAndTaxableValue{}
	eta.Municipal = model.ExemptionsAndTaxableValue{}
	eta.Independent = model.ExemptionsAndTaxableValue{}

	doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(8) > tbody > tr").Each(func(i int, s *goquery.Selection) {

		if i > 1 {
			eta = ExemptionsTaxableRecord(s, i, eta)
		}
	})

	*_bcpa.ExemptionsTaxable = eta
}

// SalesRecord Parse Sales hostory table called by LoadSalesHistory
func SalesRecord(s *goquery.Selection) model.Sale {
	sale := model.Sale{}

	s.Find("td").Each(func(int int, s *goquery.Selection) {

		switch int {
		case 0:
			sale.Date = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 1:
			sale.Type = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 2:
			sale.Price = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 3:
			sale.BookPageCIN = strings.TrimSpace(s.Find("span").First().Contents().Text())
		}

	})

	return sale
}

// LoadSalesHistory Load up the sales history table in objects and append to BCPA parent calls SalesRecord
func LoadSalesHistory(doc *goquery.Document, _bcpa *model.Bcpa) {

	doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(10) > tbody > tr > td:nth-child(1) > table:nth-child(1) > tbody > tr").Each(func(i int, s *goquery.Selection) {

		if i > 1 {
			if len(strings.TrimSpace(StripSpaces(s.Find("td:nth-child(1)").Find("span").First().Contents().Text()))) > 0 {

				sale := SalesRecord(s)
				//append the sale to the struct
				*_bcpa.SalesHistory = append(&_bcpa.SalesHistory, sale)
			}
		}
	})
}

// LandCalculationRecord extract the land calculation record called by LoadLandCalculations
func LandCalculationRecord(s *goquery.Selection) model.LandCalculation {
	lc := model.LandCalculation{}

	s.Find("td").Each(func(int int, s *goquery.Selection) {

		switch int {
		case 0:
			lc.Price = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 1:
			lc.Factor = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 2:
			lc.Type = strings.TrimSpace(s.Find("span").First().Contents().Text())
		}

	})

	return lc
}

// LoadLandCalculations load calculations structure calls LandCalculationRecord
func LoadLandCalculations(doc *goquery.Document, _bcpa *model.Bcpa) {

	//Parent node to be attached to BCPA
	lcs := model.LandCalculations{}

	//We need a Card placeholder as we'll need to set the URL for use later
	card := model.RecBuildingCard{}

	//Need to know how many rows are in the table. We only need 3-* and the last 2 or 3 rows
	rowCount := doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(10) > tbody > tr > td:nth-child(2) > table > tbody").Find("tr").Size()

	//Last row of table
	EffActYearBuiltRowIndex := rowCount - 1
	//Unit row or Bldg SF row
	UnitOrBldgRowIndex := rowCount - 2
	//Bldg Row of table
	BldgRowIndex := rowCount - 3

	//Lets loop the Table rows
	doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(10) > tbody > tr > td:nth-child(2) > table > tbody > tr").Each(func(i int, s *goquery.Selection) {

		if i == EffActYearBuiltRowIndex { //Grab the last row of the table

			lcs.EffActYearBuilt = strings.TrimSpace(StripSpaces(s.Find("td").First().Find("a").First().Find("span").First().Contents().Text()))

		} else if i == UnitOrBldgRowIndex || i == BldgRowIndex { //Grab the second to last row and check if we have Unit data or not

			//Check the value of the first td
			if strings.Contains(s.Find("td:nth-child(1)").Find("span").Contents().Text(), "Units") { //This is the unit row so grab the unit data

				lcs.Units = strings.TrimSpace(StripSpaces(s.Find("td:nth-child(2)").Find("span").Contents().Text()))

			} else { //This is the Bldg Row, Grab the building data

				//Defualt value for Units
				lcs.Units = "0"

				//Set the SF total
				lcs.AdjBldgSF = strings.TrimSpace(StripSpaces(s.Find("td:nth-child(2)").Find("span").Contents().Text()))

				var hrefExists bool

				//Grab the Sketch URL
				lcs.SketchURL, hrefExists = s.Find("td:nth-child(1)").Find("a:nth-child(3)").Attr("href")

				if !hrefExists {
					log.Println("No Sketch URL")
				}
				//Get the card URL
				card.CardURL, hrefExists = s.Find("td:nth-child(1)").Find("a:nth-child(2)").Attr("href")
				card.CardURL = url.QueryEscape(card.CardURL)

				if !hrefExists {
					log.Println("No Card URL")
				}

			}

		} else if i > 1 { //These are the data rows as we skip the header rows
			//Make sure we have data in the row before proceeding
			if len(strings.TrimSpace(StripSpaces(s.Find("td:nth-child(1)").Find("span").First().Contents().Text()))) > 0 {
				//Build the record
				LandCalculation := LandCalculationRecord(s)
				//Append the Land Cal
				lcs.Calculations = append(lcs.Calculations, LandCalculation)
			}

		}
	})

	//Add the card info if the URL isn't blank
	if card.CardURL != "" {
		lcs.Cards = append(lcs.Cards, card)
	}

	*_bcpa.LandCalculations = lcs
}

// SpecialAssessmentRecord extract data for structure called by LoadSpecialAssessments
func SpecialAssessmentRecord(s *goquery.Selection) model.SpecialAssessment {
	sa := model.SpecialAssessment{}

	s.Find("td").Each(func(int int, s *goquery.Selection) {

		switch int {
		case 0:
			sa.Fire = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 1:
			sa.Garb = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 2:
			sa.Light = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 3:
			sa.Drain = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 4:
			sa.Impr = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 5:
			sa.Safe = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 6:
			sa.Storm = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 7:
			sa.Clean = strings.TrimSpace(s.Find("span").First().Contents().Text())
		case 8:
			sa.Misc = strings.TrimSpace(s.Find("span").First().Contents().Text())
		}

	})

	return sa
}

// LoadSpecialAssessments parse assessments table calls SpecialAssessmentRecord
func LoadSpecialAssessments(doc *goquery.Document, _bcpa *model.Bcpa) {

	//Lets loop the Table rows
	doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(12) > tbody > tr").Each(func(i int, s *goquery.Selection) {
		if i > 1 {

			specialAssessment := SpecialAssessmentRecord(s)
			//append the sale to the struct
			*_bcpa.SpecialAssessments = append(&_bcpa.SpecialAssessments, specialAssessment)
		}
	})
}

//ExtractCardURL Parse the data from the card URL
func ExtractCardURL(cardURL string, i int, _bcpa *model.Bcpa) error {

	// Load the HTML document from the URL
	doc, err := goquery.NewDocument(_baseURL + cardURL)

	if err != nil {
		log.Fatal(err)
	}

	q, errParseQuery := url.Parse(cardURL)
	if err != nil {
		log.Fatal(errParseQuery)
	} else { //Since we can parse the URL lets set the values

		//urlParams := q.Query() Pulling the tax year and folio
		*_bcpa.LandCalculations.Cards[i].Folio = q.Query()["folio"][0]
		*_bcpa.LandCalculations.Cards[i].TaxYear = q.Query()["taxyear"][0]

	}

	//Grab the various values
	//Section 1
	*_bcpa.LandCalculations.Cards[i].ParcelIDNumber = SingleFindValue(doc, "#Table6 > tbody > tr:nth-child(2) > td:nth-child(1)")

	//Section 2
	*_bcpa.LandCalculations.Cards[i].UseCode = SingleFindValue(doc, "#Table7 > tbody > tr:nth-child(2) > td > p:nth-child(2) > font")

	//Section 3
	*_bcpa.LandCalculations.Cards[i].NoBedrooms = SingleFindValue(doc, "#Table1 > tbody > tr:nth-child(2) > td:nth-child(1) > p")
	_bcpa.LandCalculations.Cards[i].NoBaths = SingleFindValue(doc, "#Table1 > tbody > tr:nth-child(2) > td:nth-child(2) > p")
	**_bcpa.LandCalculations.Cards[i].NoUnits = SingleFindValue(doc, "#Table1 > tbody > tr:nth-child(2) > td:nth-child(3) > p")
	*_bcpa.LandCalculations.Cards[i].NoStories = SingleFindValue(doc, "#Table1 > tbody > tr:nth-child(2) > td:nth-child(4) > p")
	*_bcpa.LandCalculations.Cards[i].NoBuildings = SingleFindValue(doc, "#Table1 > tbody > tr:nth-child(2) > td:nth-child(5) > p")

	//Section 4
	*_bcpa.LandCalculations.Cards[i].Foundation = SingleFindValue(doc, "#Table2 > tbody > tr:nth-child(2) > td:nth-child(1) > p")
	*_bcpa.LandCalculations.Cards[i].Exterior = SingleFindValue(doc, "#Table2 > tbody > tr:nth-child(2) > td:nth-child(2) > p")
	*_bcpa.LandCalculations.Cards[i].RoofType = SingleFindValue(doc, "#Table2 > tbody > tr:nth-child(2) > td:nth-child(3) > p")
	*_bcpa.LandCalculations.Cards[i].RoofMaterial = SingleFindValue(doc, "#Table2 > tbody > tr:nth-child(2) > td:nth-child(4) > p")

	//Section 5
	*_bcpa.LandCalculations.Cards[i].Interior = SingleFindValue(doc, "#Table3 > tbody > tr:nth-child(2) > td:nth-child(1) > p")
	*_bcpa.LandCalculations.Cards[i].Floors = SingleFindValue(doc, "#Table3 > tbody > tr:nth-child(2) > td:nth-child(2) > p")
	*_bcpa.LandCalculations.Cards[i].Plumbing = SingleFindValue(doc, "#Table3 > tbody > tr:nth-child(2) > td:nth-child(3) > p")
	*_bcpa.LandCalculations.Cards[i].Electric = SingleFindValue(doc, "#Table3 > tbody > tr:nth-child(2) > td:nth-child(4) > p")
	*_bcpa.LandCalculations.Cards[i].Classification = SingleFindValue(doc, "#Table3 > tbody > tr:nth-child(2) > td:nth-child(5) > p")

	//Section 6
	*_bcpa.LandCalculations.Cards[i].CeilingHeights = SingleFindValue(doc, "#Table4 > tbody > tr:nth-child(2) > td:nth-child(1) > p")
	*_bcpa.LandCalculations.Cards[i].QualityOfConstruction = SingleFindValue(doc, "#Table4 > tbody > tr:nth-child(2) > td:nth-child(2) > p")
	*_bcpa.LandCalculations.Cards[i].CurrentConditionStructure = SingleFindValue(doc, "#Table4 > tbody > tr:nth-child(2) > td:nth-child(3) > p")
	*_bcpa.LandCalculations.Cards[i].ConstructionClass = SingleFindValue(doc, "#Table4 > tbody > tr:nth-child(2) > td:nth-child(4) > p")

	//fmt.Println(card.ParcelIDNumber)

	//Make sure we have the table
	if doc.Find("#Table8 > tbody:nth-child(1) > tr").Size() > 0 {
		fmt.Println("We Have Features")

		LoopCardFeatureTable(doc, i, &_bcpa)

	} else {
		fmt.Println("We DONT Have Features: " + strconv.Itoa(doc.Find("#Table8 > tbody:nth-child(1) > tr").Size()))
	}

	//Make sure we have permits
	if len(doc.Find("#Table5 > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1)").Find("p").Contents().Text()) > 2 {
		fmt.Println("Permits: " + strconv.Itoa(len(doc.Find("#Table5 > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1)").Find("p").Contents().Text())))
		fmt.Println("Permit 1 Val: " + doc.Find("#Table5 > tbody:nth-child(1) > tr:nth-child(2) > td:nth-child(1)").Find("p").Contents().Text())

		LoadCardPermits(doc, i, &_bcpa)

	} else {
		fmt.Println("We DONT Have Permits: " + strconv.Itoa(doc.Find("#Table5 > tbody:nth-child(1) > tr").Size()))
	}

	return err
}

// LoopCardFeatureTable parse the Features table if it exists and return a record set calls ExtractCardURL
func LoopCardFeatureTable(doc *goquery.Document, i int, _bcpa *model.Bcpa) {
	//Lets loop the Table rows

	doc.Find("#Table8 > tbody:nth-child(1) > tr").Each(func(tr int, s *goquery.Selection) {
		if tr > 1 {

			extraFeature := ExtraFeature{Feature: strings.TrimSpace(StripSpaces(s.Find("td > p").Contents().Text()))}

			//append to the struct
			*_bcpa.LandCalculations.Cards[i].ExtraFeatures = append(&_bcpa.LandCalculations.Cards[i].ExtraFeatures, extraFeature)
		}
	})
}

// LoadCardPermits load the permits from the cards page calls ExtractCardURL
func LoadCardPermits(doc *goquery.Document, i int, _bcpa *model.Bcpa) {

	permit := Permit{}

	doc.Find("#Table5 > tbody > tr").Each(func(tr int, s *goquery.Selection) {

		if tr > 1 {

			permit.PermitNo = strings.TrimSpace(StripSpaces(s.Find("td:nth-child(1)").Find("p").Contents().Text()))
			permit.PermitType = strings.TrimSpace(StripSpaces(s.Find("td:nth-child(2)").Find("p").Contents().Text()))
			permit.EstCost = strings.TrimSpace(StripSpaces(s.Find("td:nth-child(3)").Find("p").Contents().Text()))
			permit.PermitDate = strings.TrimSpace(StripSpaces(s.Find("td:nth-child(4)").Find("p").Contents().Text()))
			permit.CODate = strings.TrimSpace(StripSpaces(s.Find("td:nth-child(5)").Find("p").Contents().Text()))
			//append the permit to the struct
			*_bcpa.LandCalculations.Cards[i].Permits = append(&_bcpa.LandCalculations.Cards[i].Permits, permit)
		}
	})
}

//SingleFindValue ...
func SingleFindValue(doc *goquery.Document, exp string) string {
	return strings.TrimSpace(StripSpaces(doc.Find(exp).Contents().Text()))
}
