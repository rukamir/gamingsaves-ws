package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/render"

	_ "github.com/go-sql-driver/mysql"
)

type Game struct {
	Id    string `json:"id"`
	Title string `json:"title"`
}

func main() {
	db, err := sql.Open("mysql", "retriever:password@tcp(127.0.0.1:3306)/game")
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	log.Printf("%o", db.Stats().InUse)
	err = db.Ping()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Called")
		w.Write([]byte("welcome"))
	})
	r.Get("/v1/deals", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Called")
		var gameList []Game
		rows, err := db.Query("SELECT id, title from deals;")
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		defer rows.Close()
		for rows.Next() {
			var (
				id    string
				title string
			)
			if err := rows.Scan(&id, &title); err != nil {
				log.Fatal(err)
			}
			game := Game{Id: id, Title: title}
			gameList = append(gameList, game)
			gameSon, err := json.Marshal(game)

			if err != nil {
				panic(err.Error()) // proper error handling instead of panic in your app
			}

			log.Println(string(gameSon))
		}

		render.JSON(w, r, gameList)
	})
	http.ListenAndServe(":3000", r)
}
