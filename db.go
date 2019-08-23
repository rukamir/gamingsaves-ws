package main

import (
	"database/sql"
	"log"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// Game Data struct for handling data in DB
type Game struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Platform    string `json:"platform"`
	ListPrice   string `json:"listprice"`
	MSRP        string `json:"msrp"`
	Discount    string `json:"discount"`
	URL         string `json:"url"`
	Description string `json:"description"`
	Rating      string `json:"rating"`
	Release     string `json:"release"`
	Thumbnail   string `json:"thumbnail"`
	Source      string `json:"source"`
	Date        string `json:"date"`
}

// GameCount holds game count per platform
type GameCount struct {
	Platform string `json:"platform"`
	Count    int    `json:"count"`
}

// DB object to interact with database
var DB, err = sql.Open("mysql", "retriever:password@tcp(127.0.0.1:3306)/game")

// SetUpDB config DB
func SetUpDB() {
	DB.SetMaxOpenConns(15)
}

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
	rows, err := DB.Query("SELECT id, title, platform, list_price, msrp_price, discount, `release`, url, `source`, updated from deals;")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		var game Game
		if err := rows.Scan(&game.ID, &game.Title, &game.Platform, &game.ListPrice,
			&game.MSRP, &game.Discount, &game.Release,
			&game.URL, &game.Source, &game.Date); err != nil {
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
	var correctedVal int
	correctedVal = setWithMin(val, min)
	correctedVal = setWithMax(val, max)
	return correctedVal
}

// GetDealsQuery query db
func GetDealsQuery(order string, limit int, page int, minprice int, maxprice int, platforms string, mindiscount int) []Game {
	var gameList []Game
	limit = setWithLimits(limit, 1, 100)
	page = setWithMin(page, 1)
	mindiscount = setWithMax(mindiscount, 100)
	startIndex := (page - 1) * limit
	var queryValues []interface{}
	needsAnd := false
	completequery := "SELECT id, title, platform, list_price, msrp_price, discount, `release`, `desc`, url, thumbnail_key, rating, `source`, updated FROM deals"

	// WHERE clauses
	if minprice != 0 {
		needsAnd = true
		completequery += " WHERE list_price >= ? "
		queryValues = append(queryValues, strconv.Itoa(minprice))
	}

	if maxprice != 0 {
		if needsAnd {
			completequery += " AND "
		} else {
			completequery += " WHERE "
			needsAnd = true
		}
		completequery += " list_price <= ? "
		queryValues = append(queryValues, strconv.Itoa(maxprice))
	}

	if mindiscount != 0 {
		if needsAnd {
			completequery += " AND "
		} else {
			completequery += " WHERE "
			needsAnd = true
		}
		completequery += " discount >= ? "
		queryValues = append(queryValues, strconv.Itoa(mindiscount))
	}

	if platforms != "" {
		if needsAnd {
			completequery += " AND "
		} else {
			completequery += " WHERE "
			needsAnd = true
		}
		completequery += " platform in ("

		strplatlist := strings.Split(platforms, ",")
		var tokens []string
		var platlist []interface{}
		//tokens = strings.Repeat("?", len(platlist)).Split()
		for platform := range strplatlist {
			tokens = append(tokens, "?")
			platlist = append(platlist, platform)
		}
		completequery += strings.Join(tokens, ",") + ")"
		queryValues = append(queryValues, platlist...)
	}

	if page != 0 && limit != 0 {
		queryValues = append(queryValues, strconv.Itoa(startIndex))
		queryValues = append(queryValues, strconv.Itoa(limit))
		completequery += " LIMIT ?, ?"
	}
	completequery += ";"

	log.Printf(completequery)
	rows, err := DB.Query(completequery, queryValues...)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	defer rows.Close()
	for rows.Next() {
		var game Game
		if err := rows.Scan(&game.ID, &game.Title, &game.Platform, &game.ListPrice, &game.MSRP, &game.Discount, &game.Release, &game.Description, &game.URL, &game.Thumbnail, &game.Rating, &game.Source, &game.Date); err != nil {
			log.Fatal(err)
		}

		gameList = append(gameList, game)

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
	}
	return gameList

}
