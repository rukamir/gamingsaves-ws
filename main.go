package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

func main() {
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	SetUpDB()
	err = DB.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	r := chi.NewRouter()

	// Basic CORS
	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	r.Use(cors.Handler)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Root Called")
		w.Write([]byte("welcome"))
	})

	r.Get("/v1/deals", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Deals Called")
		order := r.URL.Query().Get("order")
		sortby := r.URL.Query().Get("sortby")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		minprice, _ := strconv.Atoi(r.URL.Query().Get("minprice"))
		maxprice, _ := strconv.Atoi(r.URL.Query().Get("maxprice"))
		mindiscount, _ := strconv.Atoi(r.URL.Query().Get("mindiscount"))
		platforms := r.URL.Query().Get("platforms")

		queryresult := GetDealsQuery(order, sortby, limit, page, minprice, maxprice, platforms, mindiscount)

		render.JSON(w, r, queryresult)
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
			"msrp", "discount", "release", "product_url", "date"})

		for _, game := range gameList {
			wcsv.Write([]string{game.ID, game.Title, game.Platform, game.ListPrice,
				game.MSRP, game.Discount, game.Release, game.URL, game.Date})
		}
		wcsv.Flush()
	})

	http.ListenAndServe(":3000", r)
}
