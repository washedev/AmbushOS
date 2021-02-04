package utils

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/mcnijman/go-emailaddress"
)

func CheckProfile(email, firstName, lastName, phoneNumber, address1, city, postcode, country string) error {

	if email == "" {
		return errors.New("Email can't be empty")
	}

	if _, err := emailaddress.Parse(email); err != nil {
		return errors.New("Invalid email")
	}

	if firstName == "" {
		return errors.New("First name can't be empty")
	}

	if lastName == "" {
		return errors.New("Last name can't be empty")
	}

	if phoneNumber == "" {
		return errors.New("Phone number can't be empty")
	}

	if address1 == "" {
		return errors.New("Address 1 can't be empty")
	}

	if city == "" {
		return errors.New("City can't be empty")
	}

	if postcode == "" {
		return errors.New("Postcode can't be empty")
	}

	if country == "" {
		return errors.New("Country can't be empty")
	}

	return nil
}

func WrapProfile(email, firstName, lastName, phoneNumber, address1 string) (string, string, string, string, string) {

	if strings.ToUpper(firstName) == "RANDOM" {
		firstName = randomdata.FirstName(randomdata.RandomGender)
	}

	if strings.ToUpper(lastName) == "RANDOM" {
		lastName = randomdata.LastName()
	}

	if strings.ToUpper(phoneNumber) == "RANDOM" {
		phoneNumber = fmt.Sprintf("%v", randomdata.Number(1000000000, 9999999999))
	}

	parsedEmail, _ := emailaddress.Parse(email)

	if strings.ToUpper(parsedEmail.LocalPart) == "RANDOM" {
		email = fmt.Sprintf("%v%v%v@%v", firstName, lastName, fmt.Sprintf("%v", randomdata.Number(100, 9999)), parsedEmail.Domain)
	}

	if strings.Contains(address1, "XXX") {
		address1 = strings.Replace(address1, "XXX", RandomString(3, "ABCDEFGHUJKLMNOPQRSTUVWXYZ"), 1)
	}

	return email, firstName, lastName, phoneNumber, address1
}
