package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DB object to interact with database
var DB, err = sql.Open("mysql", "retriever:password@tcp(127.0.0.1:3306)/game?parseTime=true")
var stmtGetTopGenreDeals *sql.Stmt
var stmtGetTopPlatformDeals *sql.Stmt
var stmtGetAllGenre *sql.Stmt
var stmtGetAllPlatforms *sql.Stmt

// SetUpDB config DB
func SetUpDB() {
	DB.SetMaxOpenConns(15)
	stmtGetTopGenreDeals, err = DB.Prepare(
		"SELECT DISTINCT " +
			"game.id, genre.title, game.platform, metacritic.score " +
			"FROM deal " +
			"LEFT JOIN game ON deal.id = game.id " +
			"INNER JOIN genre ON game.title = genre.title AND genre.genre = ? " +
			"LEFT JOIN metacritic ON game.title = metacritic.title " +
			"AND game.platform = metacritic.platform " +
			"ORDER BY " +
			"metacritic.score DESC " +
			"LIMIT ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtGetTopPlatformDeals, err = DB.Prepare(
		"SELECT game.id, game.title, game.platform, metacritic.score " +
			"FROM deal " +
			"INNER JOIN game ON deal.id = game.id AND game.platform = ? " +
			"LEFT JOIN metacritic ON game.title = metacritic.title " +
			"AND game.platform = metacritic.platform " +
			"ORDER BY " +
			"metacritic.score DESC " +
			"LIMIT ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtGetAllGenre, err = DB.Prepare("SELECT genre FROM game.genre GROUP BY genre ORDER BY genre")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtGetAllPlatforms, err = DB.Prepare("SELECT distinct game.platform FROM deal LEFT JOIN game ON deal.id = game.id")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}
func CloseDB() {
	defer stmtGetTopGenreDeals.Close()
	defer stmtGetTopPlatformDeals.Close()
}

// GetTopDealsByGenre fillout
func GetTopDealsByGenre(genre string, limit int) []GameListEntry {
	var genreDealList []GameListEntry
	var gameEntry GameListEntry
	rows, err := stmtGetTopGenreDeals.Query(genre, limit)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for rows.Next() {
		if err := rows.Scan(&gameEntry.ID, &gameEntry.Title, &gameEntry.Platform, &gameEntry.Rating); err != nil {
			log.Fatal(err)
		}
		genreDealList = append(genreDealList, gameEntry)
	}
	return genreDealList
}

// GetTopDealsByPlatform fillout
func GetTopDealsByPlatform(platform string, limit int) []GameListEntry {
	log.Printf("getting top plats")
	var genreDealList []GameListEntry
	var gameEntry GameListEntry
	rows, err := stmtGetTopPlatformDeals.Query(platform, limit)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for rows.Next() {
		if err := rows.Scan(&gameEntry.ID, &gameEntry.Title, &gameEntry.Platform, &gameEntry.Rating); err != nil {
			log.Fatal(err)
		}
		genreDealList = append(genreDealList, gameEntry)
	}
	return genreDealList
}

// GetAllPlatforms notes
func GetAllPlatforms() []string {
	var plat sql.NullString
	var platList []string
	rows, err := stmtGetAllPlatforms.Query()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for rows.Next() {
		if err := rows.Scan(&plat); err != nil {
			log.Fatal(err)
		}
		platList = append(platList, plat.String)
	}

	return platList

}

// GetAllGenres notes
func GetAllGenres() []string {
	log.Printf("getting all genres")
	var genreList []string
	var genre string
	rows, err := stmtGetAllGenre.Query()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	for rows.Next() {
		if err := rows.Scan(&genre); err != nil {
			log.Fatal(err)
		}
		genreList = append(genreList, genre)
	}

	return genreList
}
