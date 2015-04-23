package utwil

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

// Call is the Go-representation of Twilio REST API's call.
//
// Details:
//
//	https://www.twilio.com/docs/api/rest/call
//
type Call struct {
	AccountSID      string `json:"account_sid"`
	Annotation      string `json:"annotation"`
	AnsweredBy      string `json:"answered_by"`
	ApiVersion      string `json:"api_version"`
	CallerName      string `json:"caller_name"`
	DateCreated     *Time  `json:"date_created"`
	DateUpdated     *Time  `json:"date_updated"`
	Direction       string `json:"direction"`
	Duration        string `json:"duration"`
	ForwardedFrom   string `json:"forwarded_from"`
	From            string `json:"from"`
	FromFormatted   string `json:"from_formatted"`
	GroupSID        string `json:"group_sid"`
	ParentCallSID   string `json:"parent_call_sid"`
	PhoneNumberSID  string `json:"phone_number_sid"`
	Price           string `json:"price"`
	PriceUnit       string `json:"price_unit"`
	SID             string `json:"sid"`
	StartTime       *Time  `json:"start_time"`
	EndTime         *Time  `json:"end_time"`
	Status          string `json:"status"`
	SubresourceURIs struct {
		Notifications string `json:"notifications"`
		Recordings    string `json:"recordings"`
	} `json:"subresource_uris"`
	To          string `json:"to"`
	ToFormatted string `json:"to_formatted"`
	URI         string `json:"uri"`
}

// CallReq is the Go-representation of the Twilio REST API's call request.
//
// Details:
//
//	https://www.twilio.com/docs/api/rest/making-calls
//
type CallReq struct {
	From                 string
	To                   string
	URL                  string
	ApplicationSID       string
	Method               string
	FallbackURL          string
	FallbackMethod       string
	StatusCallback       string
	StatusCallbackMethod string
	SendDigits           string
	IfMachine            string
	Timeout              int
	Record               bool
}

// Submit sends a call request populating form fields only if they contain
// a non-zero value.
func (c *Client) SubmitCall(req CallReq) (*Call, error) {
	// @TODO wait until github.com/gorilla/schema supports struct-to-url.Values
	values := url.Values{}
	values.Set("From", req.From)
	values.Set("To", req.To)
	if req.URL != "" {
		values.Set("Url", req.URL)
	}
	if req.ApplicationSID != "" {
		values.Set("ApplicationSid", req.ApplicationSID)
	}
	if req.Method != "" {
		values.Set("Method", req.Method)
	}
	if req.FallbackURL != "" {
		values.Set("FallbackUrl", req.FallbackURL)
	}
	if req.FallbackMethod != "" {
		values.Set("FallbackMethod", req.FallbackMethod)
	}
	if req.StatusCallback != "" {
		values.Set("StatusCallback", req.StatusCallback)
	}
	if req.StatusCallbackMethod != "" {
		values.Set("StatusCallbackMethod", req.StatusCallbackMethod)
	}
	if req.SendDigits != "" {
		values.Set("SendDigits", req.SendDigits)
	}
	if req.IfMachine != "" {
		values.Set("IfMachine", req.IfMachine)
	}
	if req.Timeout > 0 {
		values.Set("Timeout", strconv.Itoa(req.Timeout))
	}
	if req.Record {
		values.Set("Record", "true")
	}
	call := new(Call)
	err := c.postForm(fmt.Sprintf("%s/Calls.json", c.urlPrefix()), values, call)
	return call, err
}

// Call requests the Twilio API to call a number and send a POST request to
// the given URL to report what happened:
//
// 	https://www.twilio.com/docs/api/twiml/twilio_request
//
// Example:
//
//      callbackPostURL := fmt.Sprintf(
//              "http://twimlets.com/forward?PhoneNumber=%s",
//              "+15559871234",
//      )
//      call, err := client.Call("+15551231234", "+15553214321", callbackPostURL)
//
func (c *Client) Call(from, to, callbackPostURL string) (*Call, error) {
	req := CallReq{
		From: from,
		To:   to,
		URL:  callbackPostURL,
	}
	return c.SubmitCall(req)
}

// RecordedCall is the same as Client.Call, but recorded
func (c *Client) RecordedCall(from, to, callbackPostURL string) (*Call, error) {
	req := CallReq{
		From:   from,
		To:     to,
		URL:    callbackPostURL,
		Record: true,
	}
	return c.SubmitCall(req)
}

// CallListQuery is a struct that contains an embedded utwil.ListQuery.
// The typing allows the correctly-typed iterator/list to be returned.
type CallListQuery struct{ *ListQuery }

// Calls takes a vargs of utwil.ListQueryConf functions to configure the query
// to be sent to the Twilio API:
//
// Example:
//
//	iter := client.Calls(
//		utwil.StartedBefore("2014-01-01"),
//		utwil.To("+15551231234")).Iter()
//
func (c *Client) Calls(confs ...ListQueryConf) *CallListQuery {
	return &CallListQuery{ListQuery: newListQuery(c, confs...)}
}

// StartedBefore filters calls started before a given date string "YYYY-MM-DD"
func StartedBefore(ymd string) ListQueryConf {
	return func(q *ListQuery) { q.Values.Set("StartTime<", ymd) }
}

// StartedBeforeYMD filters calls started before a given date (YMD considered only)
func StartedBeforeYMD(t time.Time) ListQueryConf {
	return StartedBefore(t.Format(YMD))
}

// StartedAfter filters calls started after a given date string "YYYY-MM-DD"
func StartedAfter(ymd string) ListQueryConf {
	return func(q *ListQuery) { q.Values.Set("StartTime>", ymd) }
}

// StartedAfterYMD filters calls started after a given date (YMD considered only)
func StartedAfterYMD(t time.Time) ListQueryConf {
	return StartedBefore(t.Format(YMD))
}

// Iter creates an iterator that iterates utwil.Call results
func (q *CallListQuery) Iter() *CallIter {
	initURI := fmt.Sprintf("%s?%s", q.callsURL(), q.Values.Encode())
	iter := &CallIter{iter: newIter(q.Client, initURI)}
	iter.iterable = new(callList)
	return iter
}

type callList struct {
	Calls []Call `json:"calls"`
	listResource
}

func (cl callList) item(idx int) interface{} { return cl.Calls[idx] }
func (cl callList) size() int                { return len(cl.Calls) }
func (cl callList) nextPage(c *Client) (iterable, error) {
	return cl.loadNextPage(c, new(callList))
}
