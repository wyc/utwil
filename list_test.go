package utwil

import (
	"log"
	"os"
	"testing"
	"time"
)

var (
	TheAccountSID      = os.Getenv("TWILIO_ACCOUNT_SID")
	TheAuthToken       = os.Getenv("TWILIO_AUTH_TOKEN")
	TheToPhoneNumber   = os.Getenv("TWILIO_DEFAULT_TO")
	TheFromPhoneNumber = os.Getenv("TWILIO_DEFAULT_FROM")
	TheClient          = New(TheAccountSID, TheAuthToken)
)

func init() {
	if TheAccountSID == "" {
		log.Fatalf("Testing env var TWILIO_ACCOUNT_SID is unset")
	} else if TheAuthToken == "" {
		log.Fatalf("Testing env var TWILIO_AUTH_TOKEN is unset")
	} else if TheToPhoneNumber == "" {
		log.Fatalf("Testing env var TWILIO_DEFAULT_TO is unset")
	} else if TheFromPhoneNumber == "" {
		log.Fatalf("Testing env var TWILIO_DEFAULT_FROM is unset")
	}
}

// Iterate (and paginate) through all the calls
func TestListCalls(t *testing.T) {
	iter := TheClient.Calls().Iter()
	callCount := 0
	var call Call
	for iter.Next(&call) {
		callCount++
	}
	if iter.Err() != nil {
		t.Fatalf("error: %s", iter.Err().Error())
	}
	t.Logf("Calls total: %d\n", callCount)
}

// Iterate (and paginate) through all calls from TheFromPhoneNumber within
// one week
func TestQueryCalls(t *testing.T) {
	weekAgo := time.Now().Add(-7 * 24 * time.Hour)
	iter := TheClient.Calls(
		From(TheFromPhoneNumber),
		StartedAfterYMD(weekAgo)).Iter()
	callCount := 0
	var call Call
	for iter.Next(&call) {
		callCount++
	}
	if iter.Err() != nil {
		t.Fatalf("error: %s", iter.Err().Error())
	}
	t.Logf("Within-one-week calls total: %d\n", callCount)
}

// Iterate (and paginate) through all the messages
func TestListMessages(t *testing.T) {
	iter := TheClient.Messages().Iter()
	msgCount := 0
	var msg Message
	for iter.Next(&msg) {
		msgCount++
	}
	if iter.Err() != nil {
		t.Fatalf("error: %s\n", iter.Err().Error())
	}
	t.Logf("Messages total: %d\n", msgCount)
}

// Iterate (and paginate) through all calls from TheFromPhoneNumber within
// one week
func TestQueryMessages(t *testing.T) {
	weekAgo := time.Now().Add(-7 * 24 * time.Hour)
	iter := TheClient.Messages(
		From(TheFromPhoneNumber),
		SentAfterYMD(weekAgo)).Iter()
	msgCount := 0
	var msg Message
	for iter.Next(&msg) {
		msgCount++
	}
	if iter.Err() != nil {
		t.Fatalf("error: %s\n", iter.Err().Error())
	}
	t.Logf("With-one-week Messages total: %d\n", msgCount)
}
