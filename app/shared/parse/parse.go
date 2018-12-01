package parse

import (
	"app/model"
	"encoding/json"
	"fmt"
	"regexp"
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

// ParseRecord PropertyAssessments  table contains the information for each user Called by LoadAppendPropertyAssessments
func ParseRecord(s *goquery.Selection) model.PropertyAssessmentValue {
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

//LoadAppendPropertyAssessments used to load and append Assessments to the BCPA parent node
func LoadAppendPropertyAssessments(doc *goquery.Document) {

	doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(6) > tbody > tr").Each(func(i int, s *goquery.Selection) {

		if i > 1 {
			pa := ParseRecord(s)
			pa.CreatedAt = time.Now()
			_bcpa.PropertyAssessments = append(_bcpa.PropertyAssessments, pa)
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
func LoadAppendExemptionsTaxable(doc *goquery.Document) {

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

	_bcpa.ExemptionsTaxable = eta
}

// SalesRecord Parse Sales hostory table called by LoadSalesHistory
func SalesRecord(s *goquery.Selection) Sale {
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
func LoadSalesHistory(doc *goquery.Document) {

	doc.Find("body > table:nth-child(3) > tbody > tr > td > table > tbody > tr:nth-child(1) > td:nth-child(1) > table:nth-child(10) > tbody > tr > td:nth-child(1) > table:nth-child(1) > tbody > tr").Each(func(i int, s *goquery.Selection) {

		if i > 1 {
			if len(strings.TrimSpace(StripSpaces(s.Find("td:nth-child(1)").Find("span").First().Contents().Text()))) > 0 {

				sale := SalesRecord(s)
				//append the sale to the struct
				_bcpa.SalesHistory = append(_bcpa.SalesHistory, sale)
			}
		}
	})
}
