package utwil

import (
	"time"
)

// YMD is used to format time.Time for querying the Twilio REST API
const YMD = "2006-01-02"

// Time is a wrapper around time.Time to support JSON marshalling to/from
// the Twilio REST API, which uses RFC1123Z.
type Time struct {
	time.Time
}

// jsonRFC1123Z is the time formatting string for utwil.Time
const jsonRFC1123Z = `"` + time.RFC1123Z + `"`

// MarshalJSON marshals time.Time into the time.RFC1123Z format
func (t *Time) MarshalJSON() ([]byte, error) {
	str := t.Format(jsonRFC1123Z)
	return []byte(str), nil
}

// UnmarshalJSON unmarshals time.Time from the time.RFC1123Z format
func (t *Time) UnmarshalJSON(data []byte) error {
	ot, err := time.Parse(jsonRFC1123Z, string(data))
	if err != nil {
		return err
	}
	t.Time = ot
	return nil
}
