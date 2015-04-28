// Package utwil contains Go utilities for dealing with the Twilio API
//
// The most-used data structure is the Client, which stores credentials and
// has useful methods for interacing with Twilio. The current supported feature
// set includes the sending of Calls and Messages, retrieval of Calls and
// Messages, and the lookup of phone numbers.
//
// These actions will incur the appropriate costs on your Twilio account.
//
// Before go test, populate env vars TWILIO_ACCOUNT_SID, TWILIO_AUTH_TOKEN,
// TWILIO_DEFAULT_TO, and TWILIO_DEFAULT_FROM.
//
// Start with:
//
//	client := utwil.NewClient(AccoutSID, AuthToken)
//
// Commonly used actions have convenience functions:
//
//	msg, err := client.SendSMS("+15551231234", "+15559879876", "Hello, world!")
//
// For more complicated requests, populate the respective XxxxxReq struct
// and call the SubmitXxxxx() method:
//
//	msgReq := utwil.MessageReq{
//        	From:           "+15559871234",
//        	To:             "+15551231234",
//        	Body:           "Hello, world!",
//        	StatusCallback: "https://post.here.com/when/msg/status/changes.twiml",
//	}
//	msg, err := client.SubmitMessage(msgReq)
//
package utwil
