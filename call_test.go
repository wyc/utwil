package utwil

import (
	"encoding/json"
	"fmt"
	"testing"
)

// This test calls ToPhoneNumber and also forwards the call to ToPhoneNumber.
// ToPhoneNumber should expect two calls.
func TestCall(t *testing.T) {
	callbackPostURL := fmt.Sprintf("http://twimlets.com/forward?PhoneNumber=%s", TheToPhoneNumber)
	call, err := TheClient.Call(TheFromPhoneNumber, TheToPhoneNumber, callbackPostURL)
	if err != nil {
		t.Fatalf("Failed: %s", err.Error())
	}
	bs, err := json.MarshalIndent(call, "", "  ")
	if err != nil {
		t.Fatalf("Failed: %s", err.Error())
	}
	t.Logf("Call:\n%s\n", string(bs))
}
