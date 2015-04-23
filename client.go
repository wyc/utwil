package utwil

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	TheBaseURL   = "https://api.twilio.com"
	TheLookupURL = "https://lookups.twilio.com/v1"

	TheAPIVersion = "2010-04-01"
)

// Client stores Twilio API credentials
type Client struct {
	accountSID string
	authToken  string
	*http.Client
}

// New creates a new instance of client, storing the AccountSID and AuthToken.
//
// These credentials can be found at:
//
//	https://www.twilio.com/user/account/settings
//
func New(accountSID, authToken string) *Client {
	if accountSID == "" {
		panic("Missing Twilio AccountSID")
	} else if authToken == "" {
		panic("Missing Twilio AuthToken")
	}

	return &Client{
		accountSID: accountSID,
		authToken:  authToken,
		Client:     &http.Client{},
	}
}

func (c *Client) getJSON(url string, result interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("GetJSON(): %s", err)
	}
	req.SetBasicAuth(c.accountSID, c.authToken)
	resp, err := c.Do(req)
	if err != nil {
		return fmt.Errorf("GetJSON(): %s", err)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("GetJSON(): %s", err)
	}

	err = decodeJSON(buf, result)
	if err != nil {
		return fmt.Errorf("GetJSON(): %s", err)
	}

	return nil
}

func (c *Client) postForm(url string, values url.Values, result interface{}) error {
	req, err := http.NewRequest("POST", url, strings.NewReader(values.Encode()))
	if err != nil {
		return fmt.Errorf("PostForm(): %s", err)
	}
	req.SetBasicAuth(c.accountSID, c.authToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.Client.Do(req)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("PostForm(): %s", err)
	}

	if result == nil {
		err = checkJSON(buf)
	} else {
		err = decodeJSON(buf, result)
	}
	if err != nil {
		return fmt.Errorf("PostForm(): %s", err)
	}

	return nil
}

func (c *Client) urlPrefix() string {
	return fmt.Sprintf("%s/%s/Accounts/%s", TheBaseURL, TheAPIVersion, c.accountSID)
}

func (c *Client) callsURL() string {
	return fmt.Sprintf("%s/Calls.json", c.urlPrefix())
}

func (c *Client) messagesURL() string {
	return fmt.Sprintf("%s/Messages.json", c.urlPrefix())
}
