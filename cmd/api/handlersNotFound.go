package main

import (
	"net/http"
	"os"
)

func (app *application) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	requestedPath := "." + r.URL.Path
	if _, err := os.Stat(requestedPath); err == nil {
		// Serve static files if they exist
		// This is necessary for serving our static swagger.yml file
		http.ServeFile(w, r, requestedPath)
	} else {
		// Return JSON response if the file does not exist
		app.notFoundResponse(w, r)
	}
}
