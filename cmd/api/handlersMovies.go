package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/calvincolton/greenlight/internal/httputils"
	"github.com/calvincolton/greenlight/internal/store"
	"github.com/calvincolton/greenlight/internal/validator"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err := httputils.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie := &store.Movie{
		Title:   input.Title,
		Year:    input.Year,
		Runtime: input.Runtime,
		Genres:  input.Genres,
	}

	v := validator.New()

	if store.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.store.Movies.Insert(movie)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/movies/%d", movie.ID))

	data := map[string]any{"movie": movie}

	err = httputils.WriteJSON(w, http.StatusCreated, data, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.ReadIDParam("movieID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.store.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	data := map[string]any{"movie": movie}

	err = httputils.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) putMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.ReadIDParam("movieID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.store.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title   string   `json:"title"`
		Year    int32    `json:"year"`
		Runtime int32    `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err = httputils.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	movie.Title = input.Title
	movie.Year = input.Year
	movie.Runtime = input.Runtime
	movie.Genres = input.Genres

	v := validator.New()

	if store.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.store.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrEditConflict):
			app.editConflictReponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	data := map[string]any{"movie": movie}

	err = httputils.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) patchMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.ReadIDParam("movieID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	movie, err := app.store.Movies.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Title   *string  `json:"title"`
		Year    *int32   `json:"year"`
		Runtime *int32   `json:"runtime"`
		Genres  []string `json:"genres"`
	}

	err = httputils.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		movie.Title = *input.Title
	}
	if input.Year != nil {
		movie.Year = *input.Year
	}
	if input.Runtime != nil {
		movie.Runtime = *input.Runtime
	}
	if input.Genres != nil {
		movie.Genres = input.Genres
	}

	v := validator.New()

	if store.ValidateMovie(v, movie); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.store.Movies.Update(movie)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrEditConflict):
			app.editConflictReponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	data := map[string]any{"movie": movie}

	err = httputils.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := httputils.ReadIDParam("movieID", r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.store.Movies.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	data := map[string]any{"message": "movie successfully deleted"}

	err = httputils.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listMoviesHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title  string
		Genres []string
		store.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.Title = httputils.ReadStrings(qs, "title", "")
	input.Genres = httputils.ReadCSV(qs, "genres", []string{})

	input.Filters.Page = httputils.ReadInt(qs, "page", 1, v)
	input.Filters.PageSize = httputils.ReadInt(qs, "page_size", 20, v)
	input.Filters.Sort = httputils.ReadStrings(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "title", "year", "runtime"}
	input.Filters.Order = httputils.ReadStrings(qs, "order", "asc")

	if store.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	movies, metadata, err := app.store.Movies.GetAll(input.Title, input.Genres, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := map[string]any{"movies": movies, "metadata": metadata}

	err = httputils.WriteJSON(w, http.StatusOK, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
