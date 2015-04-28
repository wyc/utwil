package utwil

import (
	"fmt"
	"net/url"
)

// LookupReq is the Go-representation of Twilio REST API's lookup request.
//
// Details:
//	https://www.twilio.com/docs/api/rest/lookups#lookups-query-parameters
//
type LookupReq struct {
	PhoneNumber string
	Type        string
	CountryCode string
}

// SubmitLookup sends a lookup request populating form fields only if they
// contain a non-zero value.
func (c *Client) SubmitLookup(req LookupReq) (Lookup, error) {
	// @TODO wait until github.com/gorilla/schema supports struct-to-url.Values
	values := url.Values{}
	if req.Type != "" {
		values.Add("Type", req.Type)
	}
	if req.CountryCode != "" {
		values.Add("CountryCode", req.CountryCode)
	}
	url := fmt.Sprintf("%s/PhoneNumbers/%s?%s", LookupURL, req.PhoneNumber, values.Encode())
	res := Lookup{}
	err := c.getJSON(url, &res)
	return res, err
}

// Lookup is the Go-representation of Twilio REST API's lookup.
//
// Details:
//
//      https://www.twilio.com/docs/api/rest/lookups
//
type Lookup struct {
	Carrier *struct {
		ErrorCode         *int   `json:"error_code"`
		MobileCountryCode string `json:"mobile_country_code"`
		MobileNetworkCode string `json:"mobile_network_code"`
		Name              string `json:"name"`
		Type              string `json:"type"`
	} `json:"carrier"`
	CountryCode    string `json:"country_code"`
	NationalFormat string `json:"national_format"`
	PhoneNumber    string `json:"phone_number"`
	URL            string `json:"url"`
}

// Lookup looks up a phone number's details including the carrier (Type=carrier)
//
// Example:
//
//	lookup, err := client.Lookup("+15551231234")
//      // handle err
//      fmt.Println(lookup.Carrier.Type) // "mobile", "landline", or "voip"
//
func (c *Client) Lookup(phoneNumber string) (Lookup, error) {
	req := LookupReq{
		PhoneNumber: phoneNumber,
		Type:        "carrier",
	}
	return c.SubmitLookup(req)
}

// LookupNoCarrier looks up a phone number's details without the carrier
func (c *Client) LookupNoCarrier(phoneNumber string) (Lookup, error) {
	req := LookupReq{PhoneNumber: phoneNumber}
	return c.SubmitLookup(req)
}
