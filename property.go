package wazzup

import (
	"encoding/xml"
	"fmt"

	xmldate "github.com/datainq/xml-date-time"
)

// Summary represents a Wazzup real estate property summary snapshot
type Summary struct {
	XMLName   xml.Name `xml:"RealEstatePropertySummarySnapshot"`
	ID        int      `xml:"RealEstateProperyID"`
	RealtorID int
	Updated   xmldate.CustomTime `xml:"ModificationDateTimeUtc"`
	Address   string             `xml:"AddressSummary"`
	Status    string             `xml:"RealEstateProperyStatus"`
}

// Attachment represents a real estate property's additional attachment
type Attachment struct {
	XMLName  xml.Name           `xml:"Attachment"`
	Created  xmldate.CustomTime `xml:"CreationDateTime"`
	Hash     string
	FileType string             // AVI, BMP, CSV, DOC, DOCX, FLV, GIF, JPG, MOV, MP3, MPG, MSG, PDF, PNG, PPS, PPSX, PPT, PPTX, RTF, SPD, TIF, TXT, WMV, XLS, XLSX, XML, YOUTUBE, ZIP
	Updated  xmldate.CustomTime `xml:"ModificationDateTime"`
	URL      string             `xml:"URLNormalizedFile"`
	Title    string             `xml:"Title>Translation"`
	Type     string             // PHOTO, FLOORPLAN, BROCHURE, ENERGY_CERTIFICATE, CADASTRAL_MESSAGE, CADASTRAL_MAP, VIDEO, OTHER
}

// Agency represents a property listing's agency
type Agency struct {
	XMLName xml.Name `xml:"Agency"`
	Email   string
	Logo    string `xml:"LogoURL"`
	Name    string
	Phone   string
	Website string `xml:"WebsiteURL"`
}

// Price represents a property listing's financial data
type Price struct {
	XMLName               xml.Name `xml:"Financials"`
	RentPrice             int64
	RentType              string `xml:"RentPriceType"` // PRICE_PER_MONTH, PRICE_PER_YEAR, PRICE_PER_QUARTER, PRICE_PER_HALF_YEAR, PRICE_PER_CONTRACT, M2_PER_MONTH, M2_PER_YEAR
	PurchasePrice         int64
	PurchaseCondition     string   // COSTS_BUYER, FREE_ON_NAME
	PurchaseSpecification string   // EXCLUSIVE_INTERIM_INTEREST, VAT_FISCALED, VAT_INCLUSIVE, INDEXED
	PriceCode             string   // HIGHER_BUDGET, PUBLIC_AUCTION, PRICE_IN_CONSULTATION, PRICE_TO_BE_NEGOTIATED, PRICE_ON_REQUEST, ANY_PLAUSIBLE_BID, FIXED_PRICE, BY_TENDER, ASKING_PRICE
	RentSpecification     []string `xml:"RentSpecification>Specification"` // INCL_VAT, INDEXED, INCL_SERVICE_COSTS, INCL_GAS, INCL_ELECTRICITY, INCL_WATER, INCL_FURNITURE
}

// Offer represents offer information for a property listing
type Offer struct {
	XMLName        xml.Name `xml:"Offer"`
	Acceptance     string   // IN_CONCERT, BY_DATE, DIRECT
	AcceptanceDate xmldate.CustomTime
	IsForRent      bool
	IsForSale      bool
	IsSpecial      bool
	IsTopper       bool
	StartDate      xmldate.CustomTime `xml:"AvailableFromDate"`
	EndDate        xmldate.CustomTime `xml:"AvailableUntilDate"`
}

// Location contains info on a property listing's location
type Location struct {
	XMLName    xml.Name `xml:"Address"`
	Street     string   `xml:"Streetname>Translation"`
	Number     string   `xml:"HouseNumber"`
	Addition   string   `xml:"HouseNumberPostfix"`
	PostalCode string
	District   string
	City       string  `xml:"CityName>Translation"`
	Country    string  `xml:"CountryCode"`
	Lat        float64 `xml:"Latitude"`
	Lng        float64 `xml:"Longitude"`
}

// Info contains a property's additional information
type Info struct {
	XMLName     xml.Name `xml:"PropertyInfo"`
	ID          int
	ForeignID   string
	Created     xmldate.CustomTime `xml:"CreationDateTime"`
	Updated     xmldate.CustomTime `xml:"ModificationDateTime"`
	MandateDate xmldate.CustomTime
	Status      string // AVAILABLE, SOLD_UNDER_CONDITIONS, RENTED_UNDER_CONDITIONS, SOLD, RENTED, WITHDRAWN
}

// Description represents a single property description
type Description struct {
	Language string `xml:"Language,attr"`
	Value    string `xml:",chardata"`
}

// Descriptions represents a property's descriptions
type Descriptions struct {
	Title       []*Description `xml:"Title>Translation"`
	Ad          []*Description `xml:"AdText>Translation"`
	GroundFloor []*Description `xml:"GroundFloorDescription>Translation"`
	FirstFloor  []*Description `xml:"FirstFloorDescription>Translation"`
	SecondFloor []*Description `xml:"SecondFloorDescription>Translation"`
	OtherFloor  []*Description `xml:"OtherFloorDescription>Translation"`
	Garden      []*Description `xml:"GardenDescription>Translation"`
	Balcony     []*Description `xml:"BalconyDescription>Translation"`
	Details     []*Description `xml:"DetailsDescription>Translation"`
}

// Property represents a Wazzup real estate property
type Property struct {
	XMLName      xml.Name      `xml:"RealEstateProperty"`
	Area         int64         `xml:"AreaTotals>EffectiveArea"`
	Attachments  []*Attachment `xml:"Attachments>Attachment"`
	Agency       *Agency       `xml:"Contact>Agency"`
	Bedrooms     int64         `xml:"Counts>CountOfBedrooms"`
	Rooms        int64         `xml:"Counts>CountOfRooms"`
	Descriptions *Descriptions
	Price        *Price    `xml:"Financials"`
	Location     *Location `xml:"Location>Address"`
	Offer        *Offer
	Info         *Info `xml:"PropertyInfo"`
}

// GetPropertySummary fetches a summary of properties for a realtor
func (c *Connector) GetPropertySummary(realtorID int) (*Response, error) {
	uri := fmt.Sprintf("/realestatesummary/?realtorid=%d", realtorID)
	summary, _, err := c.callGet(uri, "output")
	if err != nil {
		return nil, fmt.Errorf("could not fetch summary: %s", err)
	}

	return summary, nil
}

// GetProperty fetches a single property's details
func (c *Connector) GetProperty(realtorID int, propertyID int) (*Response, string, error) {
	uri := fmt.Sprintf("/realestate/?realtorid=%d&id=%d", realtorID, propertyID)
	property, url, err := c.callGet(uri, "output")
	if err != nil {
		return nil, "", fmt.Errorf("could not fetch property details: %s", err)
	}

	return property, url, nil
}
