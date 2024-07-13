package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	//create a router mux
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)

	mux.Get("/", app.dummyJson)

	mux.Get("/v1/dummyMovies", app.AllMovies)
	mux.Get("/v1/dynamoMovies", app.DynamoDbCreateItemDummy)
	mux.Get("/v1/localMovies", app.LocalDBCreateItem)
	mux.Get("/local/getAllMovies", app.LocalGetAllMovies)

	mux.Post("/addmovie", app.LocalPutMovie)

	// mux.Get("/v1/getvalo", app.getValoAccount)
	mux.Get("/v1/getvalo", app.getAccountDetail)

	// mux.Route("/movies", func(r chi.Router) {
	// 	r.Use(middleware.Recoverer)
	// 	r.Use(middleware.Logger)

	// 	r.Post("movies", func(w http.ResponseWriter, r *http.Request) {

	// 	})

	// 	r.Get("/", app.LocalGetAllMovies)
	// })
	return mux
}
