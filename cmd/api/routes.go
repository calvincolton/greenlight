package main

import (
	_ "embed"
	"expvar"
	"net/http"

	"github.com/calvincolton/greenlight/internal/store"
	"github.com/julienschmidt/httprouter"
)

//go:embed swagger.yaml
var swaggerYAML []byte

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

	router.NotFound = http.HandlerFunc(app.notFoundHandler)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
