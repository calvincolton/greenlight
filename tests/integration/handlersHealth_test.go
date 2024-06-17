package integration

import (
	"net/http"
	"testing"

	"github.com/calvincolton/greenlight/tests/testutils"
)

func TestHandlersHealth(t *testing.T) {
	resp, err := testutils.MakeRequest(t, http.MethodGet, "http://localhost:8081/v1/healthcheck", nil, http.StatusOK)
	if err != nil {
		t.Fatalf("could not make request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
