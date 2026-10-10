package main

import "net/http"

func (app *application) healthcheckHandler(w http.ResponseWriter, _ *http.Request) {
	if err := app.writeJSON(w, http.StatusOK, envelope{"status": "available", "environment": app.config.env, "version": version}, nil); err != nil {
		app.serverErrorResponse(w, err)
	}
}
