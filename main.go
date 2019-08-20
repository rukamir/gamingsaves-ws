package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"
)

func main() {
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	err = DB.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	r := chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Root Called")
		w.Write([]byte("welcome"))
	})

	r.Get("/v1/deals", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Deals Called")
		gameList := GetAllDeals()
		render.JSON(w, r, gameList)
	})

	r.Get("/v1/deals/count", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Count Called")
		platformList := GetPlatformCounts()
		render.JSON(w, r, platformList)
	})

	r.Get("/v1/deals/csv", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("CSV Called")
		render.JSON(w, r, "CSV unavailable")
	})

	http.ListenAndServe(":3000", r)
}
