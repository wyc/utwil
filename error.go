package utwil

import (
	"encoding/json"
	"fmt"
)

// RESTException represents an error returned by the Twilio API
//
// Details:
//
//	https://www.twilio.com/docs/errors
//
type RESTException struct {
	Code     *int        `json:"code"`
	Message  string      `json:"message"`
	MoreInfo string      `json:"more_info"`
	Status   interface{} `json:"status"`
}

// Check the returned JSON for a utwil.RESTException, and return that as an
// error if so.
func checkJSON(buf []byte) error {
	re := &RESTException{}
	err := json.Unmarshal(buf, re)
	if err != nil {
		return err
	}
	if re.Code != nil && *re.Code != 0 {
		return re
	}
	return nil
}

// Print the RESTException in a human-readable form.
func (r RESTException) Error() string {
	if r.Code != nil {
		return fmt.Sprintf("Code %d: %s", *r.Code, r.Message)
	} else if r.Status != nil {
		return fmt.Sprintf("Status %d: %s", r.Status, r.Message)
	}
	return r.Message
}
