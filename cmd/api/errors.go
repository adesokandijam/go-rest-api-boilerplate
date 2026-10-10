package main

import "net/http"

func (app *application) errorResponse(w http.ResponseWriter, status int, message any) {
	if err := app.writeJSON(w, status, envelope{"error": message}, nil); err != nil {
		app.logger.Error("response failed", "error", err)
	}
}
func (app *application) serverErrorResponse(w http.ResponseWriter, err error) {
	app.logger.Error("server error", "error", err)
	app.errorResponse(w, http.StatusInternalServerError, "the server encountered a problem")
}
func (app *application) badRequestResponse(w http.ResponseWriter, err error) {
	app.errorResponse(w, http.StatusBadRequest, err.Error())
}
func (app *application) notFoundResponse(w http.ResponseWriter, _ *http.Request) {
	app.errorResponse(w, http.StatusNotFound, "resource not found")
}
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, _ *http.Request) {
	app.errorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
}
func (app *application) failedValidationResponse(w http.ResponseWriter, errors map[string]string) {
	app.errorResponse(w, http.StatusUnprocessableEntity, errors)
}
