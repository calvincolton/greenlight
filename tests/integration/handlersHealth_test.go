package integration

import (
	"net/http"
	"testing"
	"time"

	"github.com/calvincolton/greenlight/tests/testutils"
)

func TestHandlersHealth(t *testing.T) {
	var resp *http.Response
	var err error

	url := "http://localhost:8081/v1/healthcheck"
	// Retry mechanism
	for i := 0; i < 10; i++ {
		resp = testutils.MakeRequest(t, http.MethodGet, url, nil, http.StatusOK)
		if err == nil {
			break
		}
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		t.Fatalf("could not make request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
