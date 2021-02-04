package tasks

import (
	"io/ioutil"

	"github.com/jszwec/csvutil"
)

type Row struct {
	SKU   string `csv:"SKU"`
	URL   string `csv:"URL"`
	Sizes string `csv:"SIZES"`

	Email       string `csv:"EMAIL"`
	FirstName   string `csv:"FIRST NAME"`
	LastName    string `csv:"LAST NAME"`
	PhoneNumber string `csv:"PHONE NUMBER"`
	Address1    string `csv:"ADDRESS 1"`
	Address2    string `csv:"ADDRESS 2"`
	HouseNumber string `csv:"HOUSE NUMBER"`
	State       string `csv:"STATE"`
	City        string `csv:"CITY"`
	Postcode    string `csv:"POSTCODE"`
	Country     string `csv:"COUNTRY"`
	CountryID   string `csv:"COUNTRY ID"`
	Currency    string `csv:"CURRENCY"`
}

func ReadFile(filename string) ([]Row, error) {

	fileContent, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	var rows []Row

	if err := csvutil.Unmarshal(fileContent, &rows); err != nil {
		return nil, err
	}

	return rows, nil
}
