package integration

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/calvincolton/greenlight/tests/testutils"
)

func TestShowMovieHandler(t *testing.T) {
	resp, err := testutils.MakeRequest(t, http.MethodGet, "http://localhost:8081/v1/movies/1", nil, http.StatusOK)
	if err != nil {
		t.Fatalf("could not make request: %v", err)
	}

	defer resp.Body.Close()
	var got map[string]any
	err = json.NewDecoder(resp.Body).Decode(&got)
	if err != nil {
		t.Fatalf("could not decode the response body: %v", err)
	}

	expected := map[string]any{
		"movie": map[string]any{
			"id":      1,
			"title":   "Moana",
			"year":    2018,
			"runtime": 134,
			"genres":  []string{"action", "adventure"},
			"version": 1,
		},
	}

	if !testutils.Equal(got, expected) {
		t.Errorf("expected response %v, got %v", expected, got)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, resp.StatusCode)
	}
}
