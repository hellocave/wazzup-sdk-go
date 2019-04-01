package wazzup

import (
	"encoding/xml"
	"fmt"

	xmldate "github.com/datainq/xml-date-time"
)

// Contract represents a Wazzup media contract snapshot
type Contract struct {
	XMLName    xml.Name           `xml:"MediaContractSnapshot"`
	ID         int                `xml:"MediaContractID"`
	Status     string             `xml:"MediaContractStatus"` // Active, Inactive
	Created    xmldate.CustomTime `xml:"DateTimeCreatedUtc"`
	RealtorID  int
	Name       string
	Address    string `xml:"AddressLine1"`
	PostalCode string
	City       string `xml:"CityName"`
	Region     string
	SubRegion  string
	Country    string `xml:"CountryCode"`
	Phone      string `xml:"PhoneNumber"`
	Fax        string `xml:"FaxNumber"`
	Email      string `xml:"EmailAddress"`
	Website    string `xml:"WebAddress"`
}

// IsActive checks whether a contract has an active status
func (c *Contract) IsActive() bool {
	fmt.Println(c.Status)
	return c.Status == "Active"
}

// GetContracts fetches all available contracts
func (c *Connector) GetContracts() (*Response, error) {
	contracts, err := c.callGet("/mediacontract", "activate")
	if err != nil {
		return nil, fmt.Errorf("could not fetch contracts: %s", err)
	}

	return contracts, nil
}

// ActivateContract activates a given contract
func (c *Connector) ActivateContract(contractID int) (*Response, error) {
	uri := fmt.Sprintf("/mediacontract/?id=%d&action=%s", contractID, "acceptactivation")
	r, err := c.callPost(uri, "activate", []byte{})
	if err != nil {
		return nil, fmt.Errorf("could not activate contract: %s", err)
	}

	return r, nil
}
