package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var dbAddress = os.Getenv("DB_ADDRESS")
var dbUser = os.Getenv("DB_USER")
var dbPass = os.Getenv("DB_PASS")
var dbName = os.Getenv("DB_NAME")

// DB object to interact with database
var DB, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?parseTime=true", dbUser, dbPass, dbAddress, dbName))
var stmtGetTopGenreDeals *sql.Stmt
var stmtGetTopPlatformDeals *sql.Stmt
var stmtGetTopDealsUnder *sql.Stmt
var stmtGetAllGenre *sql.Stmt
var stmtGetAllPlatforms *sql.Stmt
var stmtGetGameProfile *sql.Stmt
var stmtGetGenreByGameID *sql.Stmt
var stmtGetGenreByGameTitle *sql.Stmt
var stmtGetPriceHistLast12MonthsByID *sql.Stmt
var stmtGetGameByTitleDesc *sql.Stmt
var stmtGetGamesByMultipleGenre *sql.Stmt

// SetUpDB config DB
func SetUpDB() {
	DB.SetMaxOpenConns(15)
	stmtGetTopGenreDeals, err = DB.Prepare(
		"SELECT DISTINCT " +
			"game.id, genre.title, game.platform, metacritic.score, deal.list " +
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
		"SELECT game.id, game.title, game.platform, metacritic.score, deal.list " +
			"FROM deal " +
			"INNER JOIN game ON deal.id = game.id AND game.platform = ? " +
			"LEFT JOIN metacritic ON game.title = metacritic.title " +
			"AND game.platform = metacritic.platform " +
			"WHERE deal.list <= ? " +
			"ORDER BY " +
			"metacritic.score DESC " +
			"LIMIT ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtGetTopDealsUnder, err = DB.Prepare(
		"SELECT game.id, game.title, game.platform, metacritic.score, deal.list " +
			"FROM deal " +
			"INNER JOIN game ON deal.id = game.id AND deal.list <= ? " +
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

	stmtGetGameProfile, err = DB.Prepare("SELECT game.`id`, game.title, game.platform, `desc`, rating, `release`, msrp, current_price, pub, dev, metacritic.score, url, src " +
		"FROM game " +
		"LEFT JOIN metacritic ON game.title = metacritic.title " +
		"LEFT JOIN (SELECT `id`, `list` current_price FROM game.price_hist WHERE `id` = ? ORDER BY `date` DESC LIMIT 1) AS recent_price ON game.id = recent_price.id " +
		"WHERE game.`id` = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtGetGenreByGameTitle, err = DB.Prepare("SELECT genre from genre WHERE title = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtGetPriceHistLast12MonthsByID, err = DB.Prepare("SELECT date, list FROM game.price_hist WHERE `id` = ? AND `date` > DATE_SUB(now(), INTERVAL 12 MONTH) ORDER BY `date` DESC")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtGetGameByTitleDesc, err = DB.Prepare(
		"SELECT id, title, platform FROM game WHERE MATCH (`title`,`desc`) AGAINST (?)")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtGetGamesByMultipleGenre, err = DB.Prepare(
		"SELECT DISTINCT " +
			"game.id, genre.title, game.platform, metacritic.score " +
			"FROM " +
			"deal " +
			"LEFT JOIN game ON deal.id = game.id " +
			"RIGHT JOIN genre ON game.title = genre.title AND genre.genre in ('Action','Adventure', 'Arcade', 'Multiplayer') " +
			"LEFT JOIN metacritic ON game.title = metacritic.title " +
			"AND game.platform = metacritic.platform " +
			"GROUP BY " +
			"genre.genre, " +
			"game.title, " +
			"game.id " +
			"ORDER BY " +
			"metacritic.score DESC " +
			"LIMIT 5")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

}

// CloseDB notes
func CloseDB() {
	defer stmtGetTopGenreDeals.Close()
	defer stmtGetTopPlatformDeals.Close()
}

// GetGamesByGenreList notes
func GetGamesByGenreList(criteriaList []string) []GameListEntry {
	var games []GameListEntry
	var gameEntry GameListEntry

	// build the query string
	// https://groups.google.com/forum/#!msg/golang-nuts/vHbg09g7s2I/RKU7XsO25SIJ
	var params []interface{}
	sql := "SELECT DISTINCT " +
		"game.id, genre.title, game.platform, metacritic.score " +
		"FROM " +
		"deal " +
		"LEFT JOIN game ON deal.id = game.id " +
		"RIGHT JOIN genre ON game.title = genre.title AND genre.genre IN ( %s ) " +
		"LEFT JOIN metacritic ON game.title = metacritic.title " +
		"AND game.platform = metacritic.platform " +
		"GROUP BY " +
		"genre.genre, " +
		"game.title, " +
		"game.id " +
		"ORDER BY " +
		"metacritic.score DESC " +
		"LIMIT 10"
	var sqlIn string
	for p := range criteriaList {
		log.Printf("loop")
		params = append(params, p)
		if sqlIn != "" {
			sqlIn += ", "
		}
		sqlIn += "?"
	}
	sql = fmt.Sprintf(sql, sqlIn)
	log.Printf("%v", params)

	rows, err := DB.Query(sql, params...)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&gameEntry.ID, &gameEntry.Title, &gameEntry.Platform, &gameEntry.Score); err != nil {
			log.Fatal(err)
		}
		games = append(games, gameEntry)
	}

	return games
}

// GetGenreByTitle notes
func GetGenreByTitle(title string) []string {
	var genres []string
	var genre string
	rows, err := stmtGetGenreByGameTitle.Query(title)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&genre); err != nil {
			log.Fatal(err)
		}
		genres = append(genres, genre)
	}

	return genres
}

// GetPriceHistLast12MonthsByID notes
func GetPriceHistLast12MonthsByID(id string) []PriceHistoryDay {
	var completeHist []PriceHistoryDay
	var priceDay PriceHistoryDay
	rows, err := stmtGetPriceHistLast12MonthsByID.Query(id)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&priceDay.Date, &priceDay.ListPrice); err != nil {
			log.Fatal(err)
		}
		completeHist = append(completeHist, priceDay)
	}

	return completeHist
}

// GetGameProfile notes
func GetGameProfile(id string) GameProfile {
	var profile GameProfile
	row := stmtGetGameProfile.QueryRow(id, id)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	switch err := row.Scan(
		&profile.ID,
		&profile.Title,
		&profile.Platform,
		&profile.Desc,
		&profile.Rating,
		&profile.Release,
		&profile.MSRP,
		&profile.ListPrice,
		&profile.Publisher,
		&profile.Developer,
		&profile.Score,
		&profile.URL,
		&profile.Source); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println("worked")
	default:
		panic(err)
	}

	profile.Genres = GetGenreByTitle(profile.Title)
	profile.PriceHist = GetPriceHistLast12MonthsByID(profile.ID)

	return profile
}

// GetTopDealsByGenre fillout
func GetTopDealsByGenre(genre string, limit int) []GameListEntry {
	var genreDealList []GameListEntry
	var gameEntry GameListEntry
	rows, err := stmtGetTopGenreDeals.Query(genre, limit)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&gameEntry.ID, &gameEntry.Title, &gameEntry.Platform, &gameEntry.Score, &gameEntry.ListPrice); err != nil {
			log.Fatal(err)
		}
		genreDealList = append(genreDealList, gameEntry)
	}
	return genreDealList
}

// GetTopDealsByPlatform fillout
func GetTopDealsByPlatform(platform string, underprice int, limit int) []GameListEntry {
	var genreDealList []GameListEntry
	var gameEntry GameListEntry
	rows, err := stmtGetTopPlatformDeals.Query(platform, underprice, limit)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&gameEntry.ID, &gameEntry.Title, &gameEntry.Platform, &gameEntry.Score, &gameEntry.ListPrice); err != nil {
			log.Fatal(err)
		}
		genreDealList = append(genreDealList, gameEntry)
	}
	return genreDealList
}

// GetTopDealsUnder fillout
func GetTopDealsUnder(underprice int, limit int) []GameListEntry {
	var genreDealList []GameListEntry
	var gameEntry GameListEntry
	rows, err := stmtGetTopDealsUnder.Query(underprice, limit)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&gameEntry.ID, &gameEntry.Title, &gameEntry.Platform, &gameEntry.Score, &gameEntry.ListPrice); err != nil {
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
	defer rows.Close()
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
	var genreList []string
	var genre string
	rows, err := stmtGetAllGenre.Query()
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&genre); err != nil {
			log.Fatal(err)
		}
		genreList = append(genreList, genre)
	}

	return genreList
}

// GetGamesByTextSearch notes
func GetGamesByTextSearch(value string) []SimpleGame {
	var games []SimpleGame
	var game SimpleGame
	rows, err := stmtGetGameByTitleDesc.Query(value)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&game.ID, &game.Title, &game.Platform); err != nil {
			log.Fatal(err)
		}
		games = append(games, game)
	}

	return games

}
