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

// GameQuery holds list of games and addtional meta data
type GameQuery struct {
	Deals []Game `json:"deals"`
	Total int    `json:"total"`
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
	rows, err := DB.Query("SELECT count(*), platform FROM deals GROUP BY platform;")
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

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		gameList = append(gameList, game)
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
	correctedVal := val
	correctedVal = setWithMin(correctedVal, min)
	correctedVal = setWithMax(correctedVal, max)
	return correctedVal
}

// GetDealsQuery query db
func GetDealsQuery(order string, sortby string, limit int, page int, minprice int, maxprice int, platforms string, mindiscount int) GameQuery {
	var gameList []Game
	limit = setWithLimits(limit, 1, 120)
	page = setWithMin(page, 1)
	mindiscount = setWithMax(mindiscount, 100)
	startIndex := (page - 1) * limit
	var queryValues []interface{}
	needsAnd := false
	completequery := "SELECT id, title, platform, list_price, msrp_price, discount, `release`, `desc`, url, thumbnail_key, rating, `source`, updated FROM deals"

	// WHERE clauses
	var whereclause string
	var wherevalues []interface{}
	if minprice != 0 {
		needsAnd = true
		whereclause += " WHERE list_price >= ? "
		wherevalues = append(wherevalues, strconv.Itoa(minprice))
	}

	if maxprice != 0 {
		if needsAnd {
			whereclause += " AND "
		} else {
			whereclause += " WHERE "
			needsAnd = true
		}
		whereclause += " list_price <= ? "
		wherevalues = append(wherevalues, strconv.Itoa(maxprice))
	}

	if mindiscount != 0 {
		if needsAnd {
			whereclause += " AND "
		} else {
			whereclause += " WHERE "
			needsAnd = true
		}
		whereclause += " discount >= ? "
		wherevalues = append(wherevalues, strconv.Itoa(mindiscount))
	}

	log.Printf("Listing platforms", platforms)
	if platforms != "" {
		if needsAnd {
			whereclause += " AND "
		} else {
			whereclause += " WHERE "
			needsAnd = true
		}
		whereclause += " platform in ("

		strplatlist := strings.Split(platforms, ",")
		var tokens []string
		var platlist []interface{}

		for _, plat := range strplatlist {
			tokens = append(tokens, "?")
			platlist = append(platlist, plat)
		}
		whereclause += strings.Join(tokens, ",") + ")"
		wherevalues = append(wherevalues, platlist...)
	}
	completequery += whereclause
	queryValues = append(queryValues, wherevalues...)

	// Order
	if sortby != "" {
		var sortColumn string
		switch sortby {
		case "title":
			sortColumn = " title "
		case "price":
			sortColumn = " list_price "
		case "discount":
			sortColumn = " discount "
		case "release":
			sortColumn = " release "
		case "platform":
			sortColumn = " platform "
		default:
			sortColumn = " title "
		}

		switch order {
		case "a":
			order = " ASC "
		case "d":
			order = " DESC "
		default:
			order = " ASC "
		}
		completequery += " ORDER BY " + sortColumn + order
	}

	// Limit
	if page != 0 && limit != 0 {
		queryValues = append(queryValues, strconv.Itoa(startIndex))
		queryValues = append(queryValues, strconv.Itoa(limit))
		completequery += " LIMIT ?, ?"
	}
	completequery += ";"
	log.Printf(completequery)

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

		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}

		gameList = append(gameList, game)
	}
	// todo
	// get total count and wrap it in a object that has gameList and count
	// break out the WHERE clause for this
	var total int
	log.Print("SELECT count(*) FROM deals" + whereclause + ";")
	totalres := DB.QueryRow("SELECT count(*) FROM deals"+whereclause+";", wherevalues...)
	totalres.Scan(&total)
	return GameQuery{Deals: gameList, Total: total}
}
