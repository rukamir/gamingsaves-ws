package main

import (
	"encoding/csv"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
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

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Root Called")
		w.Write([]byte("welcome"))
	})

	r.Get("/v1/deals", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Deals Called")
		// minnn, _ := strconv.Atoi(r.URL.Query().Get("min"))
		// log.Printf("min converts to ", minnn)
		// platty := r.URL.Query().Get("min")
		// log.Printf("platforms converts to ", platty)
		order := r.URL.Query().Get("order")
		sortby := r.URL.Query().Get("sortby")
		limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
		page, _ := strconv.Atoi(r.URL.Query().Get("page"))
		minprice, _ := strconv.Atoi(r.URL.Query().Get("minprice"))
		maxprice, _ := strconv.Atoi(r.URL.Query().Get("maxprice"))
		mindiscount, _ := strconv.Atoi(r.URL.Query().Get("mindiscount"))
		platforms := r.URL.Query().Get("platforms")

		gameList := GetDealsQuery(order, sortby, limit, page, minprice, maxprice, platforms, mindiscount)

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
			"msrp", "discount", "release", "product_url", "date"})

		for _, game := range gameList {
			wcsv.Write([]string{game.ID, game.Title, game.Platform, game.ListPrice,
				game.MSRP, game.Discount, game.Release, game.URL, game.Date})
		}
		wcsv.Flush()
	})

	http.ListenAndServe(":3000", r)
}
