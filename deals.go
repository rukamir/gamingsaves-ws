package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
)

// PriceHistoryDay notes
type PriceHistoryDay struct {
	Date      time.Time `json:"date"`
	ListPrice float32   `json:"list"`
}

// GameProfile is the full details for a game
type GameProfile struct {
	ID        string            `json:"id"`
	Title     string            `json:"title"`
	Platform  NullString        `json:"platform"`
	Desc      NullString        `json:"desc"`
	Rating    NullString        `json:"rating"`
	Release   NullTime          `json:"release"`
	MSRP      float32           `json:"msrp"`
	ListPrice float32           `json:"list"`
	Score     int               `json:"score"`
	Publisher NullString        `json:"pub"`
	Developer NullString        `json:"dev"`
	Genres    []string          `json:"genres"`
	PriceHist []PriceHistoryDay `json:"history"`
	URL       string            `json:"url"`
}

// GameListEntry is a short profile of a game used to display a game in a list
type GameListEntry struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Platform  NullString `json:"platform"`
	MSRP      float32    `json:"msrp"`
	ListPrice float32    `json:"list"`
	Score     NullInt64  `json:"score"`
}

// GenreGameList hosts a genre label and the corresponding games
type GenreGameList struct {
	Genre    string          `json:"genre"`
	GameList []GameListEntry `json:"games"`
}

// All of the aliased NULLxxx code "borrowed" from :
// https://medium.com/aubergine-solutions/how-i-handled-null-possible-values-from-database-rows-in-golang-521fb0ee267

// NullInt64 is an alias for sql.NullInt64 data type
type NullInt64 struct {
	sql.NullInt64
}

// MarshalJSON for NullInt64
func (ni *NullInt64) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ni.Int64)
}

// NullBool is an alias for sql.NullBool data type
type NullBool struct {
	sql.NullBool
}

// MarshalJSON for NullBool
func (nb *NullBool) MarshalJSON() ([]byte, error) {
	if !nb.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nb.Bool)
}

// NullFloat64 is an alias for sql.NullFloat64 data type
type NullFloat64 struct {
	sql.NullFloat64
}

// MarshalJSON for NullFloat64
func (nf *NullFloat64) MarshalJSON() ([]byte, error) {
	if !nf.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(nf.Float64)
}

// NullString is an alias for sql.NullString data type
type NullString struct {
	sql.NullString
}

// MarshalJSON for NullString
func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// NullTime is an alias for mysql.NullTime data type
type NullTime struct {
	mysql.NullTime
}

// MarshalJSON for NullTime
func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		return []byte("null"), nil
	}
	val := fmt.Sprintf("\"%s\"", nt.Time.Format(time.RFC3339))
	return []byte(val), nil
}
