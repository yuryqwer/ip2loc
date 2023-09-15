package main

import "net/http"

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", app.limitRequest(http.HandlerFunc(app.home)))
	mux.Handle("/v1/report", app.setupCORS(http.HandlerFunc(app.report)))

	fileServer := http.FileServer(http.Dir("./download/"))
	mux.Handle("/v1/download/", http.StripPrefix("/v1/download", fileServer))

	return app.recoverPanic(app.logRequest(app.redirectTrailingSlash(mux)))
}
