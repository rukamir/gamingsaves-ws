package main

import (
	"encoding/csv"
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
		w.Header().Add("Content-Type", "text/csv")
		w.Header().Add("Content-disposition", "attachment; filename=dealstest.csv")

		wcsv := csv.NewWriter(w)
		gameList := GetAllDeals()

		wcsv.Write([]string{"ID", "title", "platform", "list",
			"msrp", "discount", "product_url", "date"})

		for _, game := range gameList {
			wcsv.Write([]string{game.ID, game.Title, game.Platform, game.ListPrice, game.MSRP, game.Discount, game.URL, game.Date})
		}
		wcsv.Flush()
	})

	http.ListenAndServe(":3000", r)
}
