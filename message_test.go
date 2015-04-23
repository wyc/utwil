package utwil

import (
	"encoding/json"
	"testing"
)

// This test sends a test SMS to TheToPhoneNumber
func TestSendSMS(t *testing.T) {
	msg, err := TheClient.SendSMS(TheFromPhoneNumber, TheToPhoneNumber, "Hello, world!")
	if err != nil {
		t.Fatalf("Failed: %s", err.Error())
	}
	bs, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		t.Fatalf("Failed: %s", err.Error())
	}
	t.Logf("Message Sent:\n%s\n", string(bs))

}
