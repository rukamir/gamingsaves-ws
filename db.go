package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

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
var stmtSelectDealsMostViews *sql.Stmt
var stmtSelectDealsByPlatformMostViews *sql.Stmt
var stmtSelectDealsMostRecent *sql.Stmt
var stmtSelectDealsByPlatformMostRecent *sql.Stmt
var stmtUpdateViewCountByID *sql.Stmt

// SetUpDB config DB
func SetUpDB() {
	DB.SetMaxOpenConns(15)
	stmtGetTopGenreDeals, err = DB.Prepare(
		"SELECT DISTINCT " +
			"game.id, genre.title, game.platform, metacritic.score, deal.list, game.msrp, deal.discount, game.src, game.lang, game.region " +
			"FROM deal " +
			"LEFT JOIN game ON deal.id = game.id " +
			"INNER JOIN genre ON game.title = genre.title AND genre.genre = ? " +
			"LEFT JOIN metacritic ON game.title = metacritic.title AND game.platform = metacritic.platform " +
			"WHERE game.lang = ? AND game.region = ? " +
			"ORDER BY metacritic.score DESC " +
			"LIMIT ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtGetTopPlatformDeals, err = DB.Prepare(
		"SELECT game.id, game.title, game.platform, metacritic.score, deal.list, game.msrp, deal.discount, game.src, game.lang, game.region " +
			"FROM deal " +
			"INNER JOIN game ON deal.id = game.id AND game.platform = ? " +
			"LEFT JOIN metacritic ON game.title = metacritic.title " +
			"AND game.platform = metacritic.platform " +
			"WHERE deal.list <= ? and deal.lang = ? AND deal.region = ? " +
			"ORDER BY " +
			"metacritic.score DESC " +
			"LIMIT ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtGetTopDealsUnder, err = DB.Prepare(
		"SELECT game.id, game.title, game.platform, metacritic.score, deal.list, game.msrp, deal.discount, game.src, game.lang, game.region " +
			"FROM deal " +
			"INNER JOIN game ON deal.id = game.id AND deal.list <= ? " +
			"LEFT JOIN metacritic ON game.title = metacritic.title AND game.platform = metacritic.platform " +
			"WHERE game.lang = ? AND game.region = ? " +
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

	stmtGetAllPlatforms, err = DB.Prepare("SELECT distinct game.platform FROM deal LEFT JOIN game ON deal.id = game.id WHERE game.platform IS NOT NULL")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtGetGameProfile, err = DB.Prepare(
		"SELECT game.`id`, game.title, game.platform, `desc`, rating, `release`, msrp, current_price, pub, dev, metacritic.score, url, src, lang, region " +
			"FROM game " +
			"LEFT JOIN metacritic ON game.title = metacritic.title " +
			"LEFT JOIN (SELECT `id`, `list` current_price FROM game.price_hist WHERE `id` = ? ORDER BY `date` DESC LIMIT 1) AS recent_price ON game.id = recent_price.id " +
			"WHERE game.`id` = ? AND game.lang = ? AND game.region = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtGetGenreByGameTitle, err = DB.Prepare("SELECT genre from genre WHERE title = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtGetPriceHistLast12MonthsByID, err = DB.Prepare("SELECT date, list FROM game.price_hist WHERE `id` = ? AND price_hist.lang = ? AND price_hist.region = ? AND `date` > DATE_SUB(now(), INTERVAL 12 MONTH) ORDER BY `date` DESC")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtGetGameByTitleDesc, err = DB.Prepare(
		"SELECT id, title, platform FROM game WHERE MATCH (`title`,`desc`) AGAINST (?) AND game.lang = ? AND game.region = ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// this doesnt really work. work on later?
	// and region and lang
	stmtGetGamesByMultipleGenre, err = DB.Prepare(
		"SELECT DISTINCT " +
			"game.id, genre.title, game.platform, metacritic.score " +
			"FROM " +
			"deal " +
			"LEFT JOIN game ON deal.id = game.id " +
			"RIGHT JOIN genre ON game.title = genre.title AND genre.genre in ('Action','Adventure', 'Arcade', 'Multiplayer') " +
			"LEFT JOIN metacritic ON game.title = metacritic.title AND game.platform = metacritic.platform " +
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

	stmtSelectDealsByPlatformMostViews, err = DB.Prepare(
		"SELECT game.id, game.title, game.platform, metacritic.score, deal.list, game.msrp, deal.discount, game.src " +
			"FROM deal  " +
			"INNER JOIN game ON deal.id = game.id AND game.platform = ?  " +
			"LEFT JOIN metacritic ON game.title = metacritic.title " +
			"LEFT JOIN view ON view.id = game.id AND game.platform = metacritic.platform  " +
			"WHERE game.lang = ? AND game.region = ? " +
			"ORDER BY " +
			"view.month DESC LIMIT ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtSelectDealsMostViews, err = DB.Prepare(
		"SELECT " +
			"game.id, game.title, game.platform, metacritic.score, deal.list, game.msrp, deal.discount, game.src " +
			"FROM deal  " +
			"INNER JOIN game ON deal.id = game.id " +
			"LEFT JOIN metacritic ON game.title = metacritic.title " +
			"LEFT JOIN view ON view.id = game.id AND game.platform = metacritic.platform  " +
			"WHERE game.lang = ? AND game.region = ? " +
			"ORDER BY " +
			"view.month DESC LIMIT ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtSelectDealsByPlatformMostRecent, err = DB.Prepare(
		"SELECT game.id, game.title, game.platform, metacritic.score, deal.list, game.msrp, deal.discount, game.src " +
			"FROM deal  " +
			"INNER JOIN game ON deal.id = game.id AND game.platform = ?  " +
			"LEFT JOIN metacritic ON game.title = metacritic.title " +
			"LEFT JOIN view ON view.id = game.id AND game.platform = metacritic.platform  " +
			"WHERE game.lang = ? AND game.region = ? " +
			"ORDER BY " +
			"game.release DESC LIMIT ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	stmtSelectDealsMostRecent, err = DB.Prepare(
		"SELECT " +
			"game.id, game.title, game.platform, metacritic.score, deal.list, game.msrp, deal.discount, game.src, game.lang, game.region " +
			"FROM deal  " +
			"INNER JOIN game ON deal.id = game.id " +
			"LEFT JOIN metacritic ON game.title = metacritic.title " +
			"LEFT JOIN view ON view.id = game.id AND game.platform = metacritic.platform  " +
			"WHERE game.lang = ? AND game.region = ? " +
			"ORDER BY " +
			"game.release DESC LIMIT ?")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	// todo: update this to use region and language. DB update needed too
	stmtUpdateViewCountByID, err = DB.Prepare(
		"INSERT IGNORE INTO game.view VALUES (?, 1, 1) ON DUPLICATE KEY UPDATE `month` = `month` + 1, `all` = `all` + 1")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
}

// CloseDB notes
func CloseDB() {
	defer stmtGetTopGenreDeals.Close()
	defer stmtGetTopPlatformDeals.Close()
	defer stmtGetTopDealsUnder.Close()
	defer stmtGetAllGenre.Close()
	defer stmtGetAllPlatforms.Close()
	defer stmtGetGameProfile.Close()
	defer stmtGetGenreByGameID.Close()
	defer stmtGetGenreByGameTitle.Close()
	defer stmtGetPriceHistLast12MonthsByID.Close()
	defer stmtGetGameByTitleDesc.Close()
	defer stmtGetGamesByMultipleGenre.Close()
	defer stmtSelectDealsMostViews.Close()
	defer stmtSelectDealsByPlatformMostViews.Close()
	defer stmtSelectDealsMostRecent.Close()
	defer stmtSelectDealsByPlatformMostRecent.Close()
	defer stmtUpdateViewCountByID.Close()
}

// GetGamesByGenreList notes
func GetGamesByGenreList(criteriaList []string) []GameListEntry {
	var games []GameListEntry
	var gameEntry GameListEntry

	// build the query string
	// https://groups.google.com/forum/#!msg/golang-nuts/vHbg09g7s2I/RKU7XsO25SIJ
	// todo: update for region and language
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
func GetPriceHistLast12MonthsByID(id string, lang string, region string) []PriceHistoryDay {
	var completeHist []PriceHistoryDay
	var priceDay PriceHistoryDay
	rows, err := stmtGetPriceHistLast12MonthsByID.Query(id, lang, region)
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
func GetGameProfile(id string, lang string, region string) GameProfile {
	var profile GameProfile
	row := stmtGetGameProfile.QueryRow(id, id, lang, region)
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
		&profile.Source,
		&profile.Language,
		&profile.Region); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println("worked")
	default:
		panic(err)
	}

	profile.Genres = GetGenreByTitle(profile.Title)
	profile.PriceHist = GetPriceHistLast12MonthsByID(profile.ID, profile.Language, profile.Region)

	return profile
}

// GetTopDealsByGenre fillout
func GetTopDealsByGenre(genre string, lang string, region string, limit int) []GameListEntry {
	var genreDealList []GameListEntry
	var gameEntry GameListEntry
	rows, err := stmtGetTopGenreDeals.Query(genre, lang, region, limit)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&gameEntry.ID,
			&gameEntry.Title,
			&gameEntry.Platform,
			&gameEntry.Score,
			&gameEntry.ListPrice,
			&gameEntry.MSRP,
			&gameEntry.Discount,
			&gameEntry.Source,
			&gameEntry.Language,
			&gameEntry.Region); err != nil {
			log.Fatal(err)
		}
		genreDealList = append(genreDealList, gameEntry)
	}
	return genreDealList
}

// GetTopDealsByPlatform fillout
func GetTopDealsByPlatform(platform string, underprice int, lang string, region string, limit int) []GameListEntry {
	var genreDealList []GameListEntry
	var gameEntry GameListEntry
	rows, err := stmtGetTopPlatformDeals.Query(platform, underprice, lang, region, limit)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&gameEntry.ID, &gameEntry.Title, &gameEntry.Platform, &gameEntry.Score, &gameEntry.ListPrice, &gameEntry.MSRP, &gameEntry.Discount, &gameEntry.Source, &gameEntry.Language, &gameEntry.Region); err != nil {
			log.Fatal(err)
		}
		genreDealList = append(genreDealList, gameEntry)
	}
	return genreDealList
}

// GetTopDealsUnder fillout
func GetTopDealsUnder(underprice int, lang string, region string, limit int) []GameListEntry {
	var genreDealList []GameListEntry
	var gameEntry GameListEntry
	rows, err := stmtGetTopDealsUnder.Query(underprice, lang, region, limit)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&gameEntry.ID,
			&gameEntry.Title,
			&gameEntry.Platform,
			&gameEntry.Score,
			&gameEntry.ListPrice,
			&gameEntry.MSRP,
			&gameEntry.Discount,
			&gameEntry.Source,
			&gameEntry.Language,
			&gameEntry.Region); err != nil {
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
func GetGamesByTextSearch(value string, lang string, region string) []SimpleGame {
	var games []SimpleGame
	var game SimpleGame
	rows, err := stmtGetGameByTitleDesc.Query(value, lang, region)
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

// SelectDealsMostViews note
func SelectDealsMostViews(lang string, region string, limit int) []GameListEntry {
	var games []GameListEntry
	var gameEntry GameListEntry
	rows, err := stmtSelectDealsMostViews.Query(lang, region, limit)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&gameEntry.ID, &gameEntry.Title, &gameEntry.Platform, &gameEntry.Score, &gameEntry.ListPrice, &gameEntry.MSRP, &gameEntry.Discount, &gameEntry.Source); err != nil {
			log.Fatal(err)
		}
		games = append(games, gameEntry)
	}

	return games
}

// SelectDealsByPlatformMostViews note
func SelectDealsByPlatformMostViews(platform string, lang string, region string, limit int) []GameListEntry {
	var games []GameListEntry
	var gameEntry GameListEntry
	rows, err := stmtSelectDealsByPlatformMostViews.Query(platform, lang, region, limit)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&gameEntry.ID, &gameEntry.Title, &gameEntry.Platform, &gameEntry.Score, &gameEntry.ListPrice, &gameEntry.MSRP, &gameEntry.Discount, &gameEntry.Source); err != nil {
			log.Fatal(err)
		}
		games = append(games, gameEntry)
	}

	return games
}

// SelectDealsMostRecent note
func SelectDealsMostRecent(lang string, region string, limit int) []GameListEntry {
	var games []GameListEntry
	var gameEntry GameListEntry
	rows, err := stmtSelectDealsMostRecent.Query(lang, region, limit)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(
			&gameEntry.ID,
			&gameEntry.Title,
			&gameEntry.Platform,
			&gameEntry.Score,
			&gameEntry.ListPrice,
			&gameEntry.MSRP,
			&gameEntry.Discount,
			&gameEntry.Source,
			&gameEntry.Language,
			&gameEntry.Region); err != nil {
			log.Fatal(err)
		}
		games = append(games, gameEntry)
	}

	return games
}

// SelectDealsByPlatformMostRecent note
func SelectDealsByPlatformMostRecent(platform string, lang string, region string, limit int) []GameListEntry {
	var games []GameListEntry
	var gameEntry GameListEntry
	rows, err := stmtSelectDealsByPlatformMostRecent.Query(platform, lang, region, limit)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&gameEntry.ID, &gameEntry.Title, &gameEntry.Platform, &gameEntry.Score, &gameEntry.ListPrice, &gameEntry.MSRP, &gameEntry.Discount, &gameEntry.Source); err != nil {
			log.Fatal(err)
		}
		games = append(games, gameEntry)
	}

	return games
}

// UpdateViewCountByID notes
func UpdateViewCountByID(id string) string {
	_, err = stmtUpdateViewCountByID.Exec(id)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	return ""
}

func setWithMin(val int, min int) int {
	if val >= min {
		return val
	}
	return min
}

func setWithMax(val int, max int) int {
	if val <= max {
		return val
	}
	return max
}

func setWithLimits(val int, min int, max int) int {
	correctedVal := val
	correctedVal = setWithMin(correctedVal, min)
	correctedVal = setWithMax(correctedVal, max)
	return correctedVal
}

// GetDealsQuery note
func GetDealsQuery(platform string, offset int, lang string, region string, limit int) []GameListEntry {
	var gameList []GameListEntry
	var queryValues []interface{}
	limit = setWithLimits(limit, 1, 120)
	offset = setWithMin(offset, 1)
	completequery := "SELECT game.id, game.title, game.platform, metacritic.score, deal.list, game.msrp, deal.discount, game.src, game.lang, game.region " +
		"FROM deal " +
		"INNER JOIN game ON deal.id = game.id " +
		"LEFT JOIN metacritic ON game.title = metacritic.title " +
		"AND game.platform = metacritic.platform " +
		"WHERE game.lang = ? AND game.region = ? "
	queryValues = append(queryValues, lang)
	queryValues = append(queryValues, region)

	// WHERE clauses
	var whereclause string
	var wherevalues []interface{}
	if platform != "" {
		whereclause += " AND game.platform = ? "
		wherevalues = append(wherevalues, platform)
	}

	completequery += whereclause
	queryValues = append(queryValues, wherevalues...)

	queryValues = append(queryValues, strconv.Itoa(offset))
	queryValues = append(queryValues, strconv.Itoa(limit))
	completequery += "ORDER BY title ASC LIMIT ?, ?"

	rows, err := DB.Query(completequery, queryValues...)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		var game GameListEntry

		if err := rows.Scan(
			&game.ID,
			&game.Title,
			&game.Platform,
			&game.Score,
			&game.ListPrice,
			&game.MSRP,
			&game.Discount,
			&game.Source,
			&game.Language,
			&game.Region); err != nil {
			log.Fatal(err)
		}

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		gameList = append(gameList, game)
	}
	return gameList
}
