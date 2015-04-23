package utwil

import (
	"fmt"
	"net/url"
	"time"
)

// Message is the Go-representation of Twilio REST API's message.
//
// Details:
//
//      https://www.twilio.com/docs/api/rest/message
//
type Message struct {
	AccountSID      string  `json:"account_sid"`
	APIVersion      string  `json:"api_version"`
	Body            string  `json:"body"`
	DateCreated     *Time   `json:"date_created"`
	DateSent        *Time   `json:"date_sent"`
	DateUpdated     *Time   `json:"date_updated"`
	Direction       string  `json:"direction"`
	ErrorCode       *int    `json:"error_code"`
	ErrorMessage    *string `json:"error_message"`
	From            string  `json:"from"`
	NumMedia        string  `json:"num_media"`
	NumSegments     string  `json:"num_segments"`
	Price           string  `json:"price"`
	PriceUnit       string  `json:"price_unit"`
	SID             string  `json:"sid"`
	Status          string  `json:"status"`
	SubresourceURIs struct {
		Media string `json:"media"`
	} `json:"subresource_uris"`
	To  string `json:"to"`
	URI string `json:"uri"`
}

// MessageReq is the Go-representation of Twilio REST API's message request.
//
// Details:
//
//      https://www.twilio.com/docs/api/rest/sending-messages
//
type MessageReq struct {
	From           string
	To             string
	Body           string
	MediaURL       string
	StatusCallback string
	ApplicationSID string
}

// Submit sends a message request populating form fields only if they contain
// a non-zero value.
func (c *Client) SubmitMessage(req MessageReq) (*Message, error) {
	// @TODO wait until github.com/gorilla/schema supports struct-to-url.Values
	values := url.Values{}
	values.Set("From", req.From)
	values.Set("To", req.To)
	values.Set("Body", req.Body)
	if req.MediaURL != "" {
		values.Set("MediaUrl", req.MediaURL)
	}
	if req.StatusCallback != "" {
		values.Set("StatusCallback", req.StatusCallback)
	}
	if req.ApplicationSID != "" {
		values.Set("ApplicationSid", req.ApplicationSID)
	}
	msg := new(Message)
	err := c.postForm(fmt.Sprintf("%s/Messages.json", c.urlPrefix()), values, msg)
	return msg, err
}

// SendSMS sends body from/to the specified number.
//
// Example:
//
//	msg, err := client.SendSMS("+15551231234", "+15553214321", "Hello, world!")
//
func (c *Client) SendSMS(from, to, body string) (*Message, error) {
	return c.SendMMS(from, to, body, "")
}

// SendMMS sends body and mediaURL from/to the specified number.
//
// Example:
//
// 	mediaURL := "http://i.imgur.com/sZPem77.png"
//      body := "Hello, world!"
//      msg, err := client.SendMMS("+15551231234", "+15553214321", body, mediaURL)
//
func (c *Client) SendMMS(from, to, body, mediaURL string) (*Message, error) {
	req := MessageReq{
		From:     from,
		To:       to,
		Body:     body,
		MediaURL: mediaURL,
	}
	return c.SubmitMessage(req)
}

// MessageListQuery is a struct that contains an embedded utwil.ListQuery.
// The typing allows the correctly-typed iterator/list to be returned.
type MessageListQuery struct{ *ListQuery }

// Messages takes a vargs of utwil.ListQueryConf functions to configure the
// query to be sent to the Twilio API:
//
//      iter := client.Messages(
//              utwil.SentAfter("2014-01-01"),
//              utwil.From("+15551231234")).Iter()
//
func (c *Client) Messages(confs ...ListQueryConf) *MessageListQuery {
	return &MessageListQuery{ListQuery: newListQuery(c, confs...)}
}

// SentBefore filters messages sent before a given date string "YYYY-MM-DD"
func SentBefore(ymd string) ListQueryConf {
	return func(q *ListQuery) { q.Values.Set("DateSent<", ymd) }
}

// SentBeforeYMD filters messages sent before a given date (YMD considered only)
func SentBeforeYMD(t time.Time) ListQueryConf {
	return SentBefore(t.Format(YMD))
}

// SentAfter filters messages sent after a given date string "YYYY-MM-DD"
func SentAfter(ymd string) ListQueryConf {
	return func(q *ListQuery) { q.Values.Set("DateSent>", ymd) }
}

// SentAfterYMD filters messages sent after a given date (YMD considered only)
func SentAfterYMD(t time.Time) ListQueryConf {
	return SentAfter(t.Format(YMD))
}

// Iter creates an iterator that iterates utwil.Message results
func (q *MessageListQuery) Iter() *MessageIter {
	initURI := fmt.Sprintf("%s?%s", q.messagesURL(), q.Values.Encode())
	iter := &MessageIter{iter: newIter(q.Client, initURI)}
	iter.iterable = new(messageList)
	return iter
}

type messageList struct {
	Messages []Message `json:"messages"`
	listResource
}

func (ml messageList) item(idx int) interface{} { return ml.Messages[idx] }
func (ml messageList) size() int                { return len(ml.Messages) }
func (ml messageList) nextPage(c *Client) (iterable, error) {
	return ml.loadNextPage(c, new(messageList))
}
