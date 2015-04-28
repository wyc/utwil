package utwil

import (
	"encoding/json"
	"testing"
)

// This test sends a test SMS to ToPhoneNumber
func TestSendSMS(t *testing.T) {
	msg, err := TestClient.SendSMS(FromPhoneNumber, ToPhoneNumber, "Hello, world!")
	if err != nil {
		t.Fatalf("Failed: %s", err.Error())
	}
	bs, err := json.MarshalIndent(msg, "", "  ")
	if err != nil {
		t.Fatalf("Failed: %s", err.Error())
	}
	t.Logf("Message Sent:\n%s\n", string(bs))

}
