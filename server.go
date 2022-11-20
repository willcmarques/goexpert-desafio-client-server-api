package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
)

const dolar_price_url = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

type USDBRL struct {
	Bid          string `json:"bid"`
	ExchangeDate string `json:"create_date"`
}

type AwesomeAPIResponse struct {
	USDBRL USDBRL `json:"USDBRL"`
}

type DollarExchangeResponse struct {
	DollarPrice float64 `json:"dollar_price"`
}

func main() {
	setupDatabase()
	mux := http.NewServeMux()
	mux.HandleFunc("/cotacao", cotacaoHandler)
	http.ListenAndServe(":8080", mux)
}

func setupDatabase() {
	db := getDBConnection()
	_, err := db.Exec("CREATE TABLE IF NOT EXISTS exchange_rate(id TEXT, dollar_price REAL, created_at TEXT)")
	if err != nil {
		log.Fatalf("Error when creating table: %v", err)
	}
	db.Close()
}

func getDBConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./database/dollar.db")
	if err != nil {
		log.Fatalf("Error when connecting to the database: %v", err)
	}
	return db
}

func cotacaoHandler(w http.ResponseWriter, r *http.Request) {
	exchange, err := getDollarExchange()
	if err != nil {
		handlerError(w, fmt.Sprintf("Error when retrieving dollar price: %v", err))
		return
	}

	var apiResponse AwesomeAPIResponse
	err = json.Unmarshal([]byte(exchange), &apiResponse)
	if err != nil {
		handlerError(w, fmt.Sprintf("Error when unmarshalling response from API: %v", err))
		return
	}

	dollarPrice, err := strconv.ParseFloat(apiResponse.USDBRL.Bid, 64)
	if err != nil {
		handlerError(w, fmt.Sprintf("Error when converting the value %v: %v", apiResponse.USDBRL.Bid, err))
		return
	}

	err = saveDollarPrice(dollarPrice)
	if err != nil {
		handlerError(w, fmt.Sprintf("Error when try to save in database: %v", err))
		return
	}

	json.NewEncoder(w).Encode(&DollarExchangeResponse{DollarPrice: dollarPrice})

}

func handlerError(w http.ResponseWriter, message string) {
	log.Print(message)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}

func getDollarExchange() (string, error) {
	res, err := http.Get(dolar_price_url)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Error when calling API: %v", err))
	}
	defer res.Body.Close()
	content, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf(fmt.Sprintf("Response conversion error: %v", err))
	}
	return string(content), nil
}

func saveDollarPrice(dollarPrice float64) error {
	db := getDBConnection()
	defer db.Close()
	_, err := db.Exec("INSERT INTO exchange_rate(id, dollar_price, created_at) VALUES(?,?,?)", uuid.New().String(), dollarPrice, time.Now())
	if err != nil {
		return err
	}
	return nil
}
