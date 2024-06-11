package main

import (
	"fmt"
	"reflect"
)

func newTestApplication() *application {
	app := &application{}
	app.config.cors.trustedOrigins = []string{"https://trusted.com"}

	return app
}

func envelopeEqual(got, expected envelope) bool {
	for key, expectedValue := range expected {
		gotValue, exists := got[key]
		if !exists {
			printMismatch(key, gotValue, expectedValue, "Key not found")
			return false
		}

		if !reflect.DeepEqual(gotValue, expectedValue) {
			printMismatch(key, gotValue, expectedValue, "Value mismatch")
			return false
		}
	}

	for key, gotValue := range got {
		expectedValue, exists := expected[key]
		if !exists {
			printMismatch(key, gotValue, expectedValue, "Unexpected key")
			return false
		}

		if !reflect.DeepEqual(gotValue, expectedValue) {
			printMismatch(key, gotValue, expectedValue, "Value mismatch")
			return false
		}
	}

	return true
}

func printMismatch(key string, gotValue, expectedValue any, reason string) {
	fmt.Printf("Mismatch at key: %s\nReason: %s\nGot value: %#v\nExpected value: %#v\n\n", key, reason, gotValue, expectedValue)
}
