package main

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/render"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
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
	defer CloseDB()

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
		GetTopDealsByGenre("Action", 5)
		w.Write([]byte("welcome"))
	})

	r.Get("/game/{id}", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Game called")
		render.JSON(w, r, GetGameProfile(chi.URLParam(r, "id")))
	})

	r.Get("/simple-search", func(w http.ResponseWriter, r *http.Request) {
		term := r.URL.Query().Get("value")
		list := GetGamesByTextSearch(term)
		// render.JSON(w, r, GetGameProfile(chi.URLParam(r, "id")))
		render.JSON(w, r, list)
	})

	r.Get("/top/genre/multi", func(w http.ResponseWriter, r *http.Request) {
		genresRaw := r.URL.Query().Get("values")
		searchList := strings.Split(genresRaw, ",")
		for _, v := range searchList {
			log.Printf("%s", v)
		}
		list := GetGamesByGenreList(searchList)
		// render.JSON(w, r, GetGameProfile(chi.URLParam(r, "id")))
		render.JSON(w, r, list)
	})

	r.Get("/top/genre/picks", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Genre Called")
		genreList := []string{"Action", "Adventure", "Arcade", "Fighting", "First-Person", "Indie", "Platformer", "Racing", "Role-Playing", "Sports", "Strategy", "Puzzle"}
		var topGamesPerGenreList []CategoryGameList
		var gListEntry CategoryGameList

		for _, val := range genreList {
			gListEntry.Category = val
			gListEntry.GameList = GetTopDealsByGenre(val, 10)
			topGamesPerGenreList = append(topGamesPerGenreList, gListEntry)
		}

		render.JSON(w, r, topGamesPerGenreList)
	})

	r.Get("/top/genre/all", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Genre Called")
		genreList := GetAllGenres()
		var topGamesPerGenreList []CategoryGameList
		var gListEntry CategoryGameList

		for _, val := range genreList {
			gListEntry.Category = val
			gListEntry.GameList = GetTopDealsByGenre(val, 10)
			topGamesPerGenreList = append(topGamesPerGenreList, gListEntry)
		}

		render.JSON(w, r, topGamesPerGenreList)
	})

	r.Get("/top/platform", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Top Games By Platform Under $$")
		val := r.URL.Query().Get("value")
		listunder, _ := strconv.Atoi(r.URL.Query().Get("listunder"))
		if listunder == 0 {
			listunder = 1000
		}

		render.JSON(w, r, GetTopDealsByPlatform(val, listunder, 10))
	})

	r.Get("/top", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Top Games Under $$")
		listunder, _ := strconv.Atoi(r.URL.Query().Get("under"))
		if listunder == 0 {
			listunder = 1000
		}

		render.JSON(w, r, GetTopDealsUnder(listunder, 10))
	})

	r.Get("/top/platform/modern", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Top Platforms called")
		var topGamesPerPlatform []CategoryGameList
		var platEntry CategoryGameList
		platList := []string{"Nintendo Switch", "Xbox One", "PS4"}
		for _, val := range platList {
			platEntry.Category = val
			platEntry.GameList = GetTopDealsByPlatform(val, 10000, 10)
			topGamesPerPlatform = append(topGamesPerPlatform, platEntry)
		}

		render.JSON(w, r, topGamesPerPlatform)
	})

	r.Get("/top/platform/all", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Top Platforms called")
		var topGamesPerPlatform []CategoryGameList
		var platEntry CategoryGameList
		platList := GetAllPlatforms()
		for _, val := range platList {
			platEntry.Category = val
			platEntry.GameList = GetTopDealsByPlatform(val, 1000, 10)
			topGamesPerPlatform = append(topGamesPerPlatform, platEntry)
		}

		render.JSON(w, r, topGamesPerPlatform)
	})

	r.Get("/popular", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Platforms Available called")
		platList := GetAllPlatforms()

		render.JSON(w, r, platList)
	})

	r.Get("/v1/platform/available", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Platforms Available called")
		platList := GetAllPlatforms()

		render.JSON(w, r, platList)
	})

	// r.Get("/v1/deals", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Printf("Deals Called")
	// 	order := r.URL.Query().Get("order")
	// 	sortby := r.URL.Query().Get("sortby")
	// 	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	// 	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	// 	minprice, _ := strconv.Atoi(r.URL.Query().Get("minprice"))
	// 	maxprice, _ := strconv.Atoi(r.URL.Query().Get("maxprice"))
	// 	mindiscount, _ := strconv.Atoi(r.URL.Query().Get("mindiscount"))
	// 	platforms := r.URL.Query().Get("platforms")

	// 	queryresult := GetDealsQuery(order, sortby, limit, page, minprice, maxprice, platforms, mindiscount)

	// 	render.JSON(w, r, queryresult)
	// })

	// r.Get("/v1/deals/count", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Printf("Count Called")
	// 	platformList := GetPlatformCounts()
	// 	render.JSON(w, r, platformList)
	// })

	// r.Get("/v1/deals/csv", func(w http.ResponseWriter, r *http.Request) {
	// 	log.Printf("CSV Called")
	// 	w.Header().Add("Content-Type", "text/csv")
	// 	w.Header().Add("Content-disposition", "attachment; filename=dealstest.csv")

	// 	wcsv := csv.NewWriter(w)
	// 	gameList := GetAllDeals()

	// 	wcsv.Write([]string{"ID", "title", "platform", "list",
	// 		"msrp", "discount", "release", "product_url", "date"})

	// 	for _, game := range gameList {
	// 		wcsv.Write([]string{game.ID, game.Title, game.Platform, game.ListPrice,
	// 			game.MSRP, game.Discount, game.Release, game.URL, game.Date})
	// 	}
	// 	wcsv.Flush()
	// })

	log.Printf("Starting server")
	http.ListenAndServe(":2000", r)
}
