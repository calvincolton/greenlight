package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheckHandler(t *testing.T) {
	app := &application{
		config: config{
			env: "development",
		},
	}

	req, err := http.NewRequest(http.MethodGet, "/v1/healthcheck", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	app.healthcheckHandler(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", rr.Code, http.StatusOK)
	}

	expected := envelope{
		"status": "available",
		"system_info": map[string]any{
			"environment": "development",
			"version":     version,
		},
	}

	var got envelope
	err = json.NewDecoder(rr.Body).Decode(&got)
	if err != nil {
		t.Fatal(err)
	}

	if !envelopeEqual(got, expected) {
		t.Errorf("handler returned unexpected body: got %v want %v", got, expected)
	}
}
