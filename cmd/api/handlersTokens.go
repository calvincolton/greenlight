package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/calvincolton/greenlight/internal/httputils"
	"github.com/calvincolton/greenlight/internal/store"
	"github.com/calvincolton/greenlight/internal/validator"
)

func (app *application) createAuthenticaitonTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := httputils.ReadJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	store.ValidateEmail(v, input.Email)
	store.ValidatePasswordPlaintext(v, input.Password)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.store.Users.GetByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrRecordNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	token, err := app.store.Tokens.New(user.ID, 24*time.Hour, store.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := map[string]any{"authentication_token": token}

	err = httputils.WriteJSON(w, http.StatusCreated, data, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
