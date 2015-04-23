package utwil

import (
	"fmt"
	"net/url"
	"reflect"
	"sync"
)

// ListQuery stores query filter configuration and a *utwil.Client,
// and is used by MessageListQuery and CallListQuery.
type ListQuery struct {
	url.Values
	*Client
}

// ListQueryConf configures a passed *utwil.ListQuery
type ListQueryConf func(*ListQuery)

// newListQuery takes a utwil.Client and functional options to create and
// configure a new ListQuery.
//
// Read more about "functional options":
//
//      http://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
//
func newListQuery(c *Client, confs ...ListQueryConf) *ListQuery {
	q := &ListQuery{
		Values: make(url.Values),
		Client: c,
	}
	for _, conf := range confs {
		conf(q)
	}
	return q

}

// From filters calls and messages sent from a phone number.
func From(phoneNumber string) ListQueryConf {
	return func(q *ListQuery) { q.Values.Set("From", phoneNumber) }
}

// To filters calls and messages sent to a phone number.
func To(phoneNumber string) ListQueryConf {
	return func(q *ListQuery) { q.Values.Set("To", phoneNumber) }
}

type listResource struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	NumPages int `json:"num_pages"`

	Start int `json:"start"`
	End   int `json:"end"`
	Total int `json:"total"`

	URI             string  `json:"uri"`
	PreviousPageURI *string `json:"previous_page_uri"`
	NextPageURI     *string `json:"next_page_uri"`
	FirstPageURI    string  `json:"first_page_uri"`
	LastPageURI     string  `json:"last_page_uri"`
}

func (lr listResource) nextPageFullURI() string {
	return fmt.Sprintf("%s%s", TheBaseURL, *lr.NextPageURI)
}

func (lr listResource) loadNextPage(c *Client, result iterable) (iterable, error) {
	err := c.getJSON(lr.nextPageFullURI(), result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (lr listResource) hasNextPage() bool {
	return lr.NextPageURI != nil && *lr.NextPageURI != ""
}

type iterable interface {
	item(idx int) interface{}
	size() int

	hasNextPage() bool
	nextPage(*Client) (iterable, error)
}

type iter struct {
	m        sync.Mutex
	err      error
	iterable iterable
	pageItem int
	didInit  bool
	initURI  string
	client   *Client
}

func newIter(c *Client, initURI string) *iter {
	return &iter{
		m:        sync.Mutex{},
		err:      nil,
		iterable: nil,
		pageItem: 0,
		didInit:  false,
		initURI:  initURI,
		client:   c,
	}
}

func (iter *iter) loadInitURI() error {
	if iter.iterable == nil {
		panic("iterable uninitalized")
	} else if iter.initURI == "" {
		return fmt.Errorf("initURI uninitialized")
	}

	err := iter.client.getJSON(iter.initURI, iter.iterable)
	if err != nil {
		return err
	}
	return nil
}

func (iter *iter) next(result interface{}) bool {
	iter.m.Lock()
	defer iter.m.Unlock()

	if !iter.didInit {
		err := iter.loadInitURI()
		if err != nil {
			iter.err = err
			return false
		}
		iter.pageItem = 0
		iter.didInit = true
	}

	if iter.pageItem == iter.iterable.size() {
		if !iter.iterable.hasNextPage() {
			return false
		}
		nextIter, err := iter.iterable.nextPage(iter.client)
		if err != nil {
			iter.err = err
			return false
		}
		iter.pageItem = 0
		iter.iterable = nextIter
	}

	item := iter.iterable.item(iter.pageItem)
	lValuePtr := reflect.ValueOf(result)
	lValue := reflect.Indirect(lValuePtr)
	rValue := reflect.ValueOf(item)
	if lValue.Type() != rValue.Type() {
		panic(fmt.Sprintf("Iter.next() tried to load %s into %s",
			rValue.Type(), lValue.Type()))
	}
	lValue.Set(rValue)
	iter.pageItem++
	return true
}

// Err returns the latest error a MessageIter or CallIter encounters or nil
// if there was no error.
func (iter *iter) Err() error { return iter.err }

// MessageIter iterates through Twilio messages.
type MessageIter struct{ *iter }

// Next attempts to populate msg with the next utwil.Message, returning false
// if it could not due to out of messages or an error. It is therefore
// recommended to check for errors with MessageIter.Err() after use:
//
// Example:
//
//	var msg utwil.Message
// 	for messageIter.Next(&msg) {
//		// use msg
//	}
//	if messageIter.Err() != nil {
//		// handle err
//	}
//
func (iter *MessageIter) Next(msg *Message) bool { return iter.next(msg) }

// CallIter iterates through Twilio calls.
type CallIter struct{ *iter }

// Next attempts to populate call with the next utwil.Call, returning false
// if it could not due to out of messages or an error. It is therefore
// recommended to check for errors with CallIter.Err() after use:
//
// Example:
//
//	var call utwil.Call
// 	for callIter.Next(&call) {
//		// use msg
//	}
//	if callIter.Err() != nil {
//		// handle err
//	}
//
func (iter *CallIter) Next(call *Call) bool { return iter.next(call) }
