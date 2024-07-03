package main

import (
	"expvar"
	"net/http"

	"github.com/calvincolton/greenlight/internal/store"
	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// OpenAPI/swagger
	// Serve the Swagger UI files
	router.ServeFiles("/swagger-ui/*filepath", http.Dir("swagger-ui/dist"))
	router.GET("/swagger.yaml", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.ServeFile(w, r, "swagger.yaml")
	})
	router.GET("/docs", func(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
		http.Redirect(w, r, "/swagger-ui/index.html?url=/swagger.yaml", http.StatusFound)
	})

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

	// TODO: hide route via load balancer / reverse proxy so only available locally
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}
