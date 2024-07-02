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
			"year":    2016,
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

func TestListMoviesHandler(t *testing.T) {
	resp, err := testutils.MakeRequest(t, http.MethodGet, "http://localhost:8081/v1/movies", nil, http.StatusOK)
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
		"movies": []map[string]any{
			{
				"id":      1,
				"title":   "Moana",
				"year":    2016,
				"runtime": 134,
				"genres":  []string{"action", "adventure"},
				"version": 1,
			},
			{
				"id":      2,
				"title":   "Deadpool",
				"year":    2016,
				"runtime": 108,
				"genres":  []string{"action", "comedy", "superhero"},
				"version": 1,
			},
			{
				"id":      3,
				"title":   "The Breakfast Club",
				"year":    1986,
				"runtime": 96,
				"genres":  []string{"drama"},
				"version": 1,
			},
		},
		"metadata": map[string]any{
			"current_page":  1,
			"page_size":     20,
			"first_page":    1,
			"last_page":     1,
			"total_records": 3,
		},
	}

	if !testutils.Equal(got, expected) {
		t.Errorf("expected response: %v, got %v", expected, got)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status %d, got %d", resp.StatusCode, http.StatusOK)
	}
}
