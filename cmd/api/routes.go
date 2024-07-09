package main

import (
	"expvar"
	"net/http"
	"path/filepath"

	"github.com/calvincolton/greenlight/internal/store"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	// healthcheck
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	// movies
	router.HandlerFunc(http.MethodPost, "/v1/movies", app.requirePermission(store.PermissionMoviesWrite, app.createMovieHandler))
	router.HandlerFunc(http.MethodGet, "/v1/movies", app.listMoviesHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movies/:movieID", app.showMovieHandler)
	// if you want to have resources accessible only to authenticated users:
	// router.HandlerFunc(http.MethodGet, "/v1/movies", app.requirePermission(store.PermissionMoviesRead, app.listMoviesHandler))
	// router.HandlerFunc(http.MethodGet, "/v1/movies/:movieID", app.requirePermission(store.PermissionMoviesRead, app.showMovieHandler))
	router.HandlerFunc(http.MethodPut, "/v1/movies/:movieID", app.requirePermission(store.PermissionMoviesWrite, app.putMovieHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/movies/:movieID", app.requirePermission(store.PermissionMoviesWrite, app.patchMovieHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/movies/:movieID", app.requirePermission(store.PermissionMoviesWrite, app.deleteMovieHandler))

	// users
	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	// authentication
	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticaitonTokenHandler)

	// metrics
	// TODO: hide route via load balancer / reverse proxy so only available locally
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	// Swagger UI and OpenAPI spec
	swaggerUIDir := http.Dir(filepath.Join(".", "swagger-ui", "dist"))
	router.Handler(http.MethodGet, "/v1/docs/*filepath", http.StripPrefix("/v1/docs/", http.FileServer(swaggerUIDir)))

	// Serve the OpenAPI spec file
	apiSpecDir := http.Dir(filepath.Join("cmd", "api"))
	router.Handler(http.MethodGet, "/v1/swagger-v1.yaml", http.StripPrefix("/v1/", http.FileServer(apiSpecDir)))

	router.NotFound = http.HandlerFunc(app.notFoundHandler)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
