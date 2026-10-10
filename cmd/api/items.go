package main

import (
	"errors"
	"github.com/adesokandijam/go-rest-api-boilerplate/internal/data"
	"github.com/adesokandijam/go-rest-api-boilerplate/internal/validator"
	"net/http"
)

func (app *application) createItemHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, err)
		return
	}
	item := &data.Item{Name: input.Name, Description: input.Description}
	v := validator.New()
	data.ValidateItem(v, item)
	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}
	if err := app.models.Items.Insert(item); err != nil {
		app.serverErrorResponse(w, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", "/v1/items/"+stringID(item.ID))
	if err := app.writeJSON(w, http.StatusCreated, envelope{"item": item}, headers); err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app *application) showItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	item, err := app.models.Items.Get(id)
	if errors.Is(err, data.ErrRecordNotFound) {
		app.notFoundResponse(w, r)
		return
	}
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"item": item}, nil); err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app *application) listItemsHandler(w http.ResponseWriter, _ *http.Request) {
	items, err := app.models.Items.GetAll()
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"items": items}, nil); err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app *application) updateItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	item, err := app.models.Items.Get(id)
	if errors.Is(err, data.ErrRecordNotFound) {
		app.notFoundResponse(w, r)
		return
	}
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}
	var input struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
	if err := app.readJSON(w, r, &input); err != nil {
		app.badRequestResponse(w, err)
		return
	}
	if input.Name != nil {
		item.Name = *input.Name
	}
	if input.Description != nil {
		item.Description = *input.Description
	}
	v := validator.New()
	data.ValidateItem(v, item)
	if !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}
	if err := app.models.Items.Update(item); err != nil {
		app.serverErrorResponse(w, err)
		return
	}
	if err := app.writeJSON(w, http.StatusOK, envelope{"item": item}, nil); err != nil {
		app.serverErrorResponse(w, err)
	}
}

func (app *application) deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	err = app.models.Items.Delete(id)
	if errors.Is(err, data.ErrRecordNotFound) {
		app.notFoundResponse(w, r)
		return
	}
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
