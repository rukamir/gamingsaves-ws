package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// Game Data struct for handling data in DB
type Game struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Platform  string `json:"platform"`
	ListPrice string `json:"listprice"`
	MSRP      string `json:"msrp"`
	Discount  string `json:"discount"`
	URL       string `json:"url"`
	Source    string `json:"source"`
	Date      string `json:"date"`
}

// GameCount holds game count per platform
type GameCount struct {
	Platform string `json:"platform"`
	Count    int    `json:"count"`
}

// DB object to interact with database
var DB, err = sql.Open("mysql", "retriever:password@tcp(127.0.0.1:3306)/game")

// GetPlatformCounts returns a list of numnber of games per platform
func GetPlatformCounts() []GameCount {
	rows, err := DB.Query("SELECT count(*), platform from deals;")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	var platcountlist []GameCount
	for rows.Next() {
		var platcount GameCount
		if err := rows.Scan(&platcount.Count, &platcount.Platform); err != nil {
			log.Fatal(err)
		}
		platcountlist = append(platcountlist, platcount)
	}
	return platcountlist
}

// GetAllDeals returns [some struct] for all REST data
func GetAllDeals() []Game {
	var gameList []Game
	rows, err := DB.Query("SELECT id, title, platform, list_price, msrp_price, discount, url, source, created from deals;")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		var game Game
		if err := rows.Scan(&game.ID, &game.Title, &game.Platform, &game.ListPrice, &game.MSRP, &game.Discount, &game.URL, &game.Source, &game.Date); err != nil {
			log.Fatal(err)
		}
		// game := Game{ID: id, Title: title}
		gameList = append(gameList, game)

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
	return gameList
}
