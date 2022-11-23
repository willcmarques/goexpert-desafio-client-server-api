package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

const DOLAR_PRICE_URL = "http://localhost:8080/cotacao"

type DollarExchangeResponse struct {
	DollarPrice float64 `json:"dollar_price"`
}

func main() {
	dollarPrice, err := getDollarPrice()
	if err != nil {
		log.Fatalln(err)
	}
	err = save(dollarPrice)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Dollar price saved successfully!")
}

func save(dollarPrice float64) error {
	file, err := os.OpenFile("client/cotacao.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("error when creating file: %v", err)
	}
	_, err = file.WriteString(fmt.Sprintf("Dólar: %v\n", dollarPrice))
	if err != nil {
		return fmt.Errorf("error when writing on file: %v", err)
	}
	return nil
}

func getDollarPrice() (float64, error) {
	c := http.Client{Timeout: time.Millisecond * 300}
	res, err := c.Get(DOLAR_PRICE_URL)
	if err != nil {
		return 0, fmt.Errorf("error on API call: %v", err)
	}
	defer res.Body.Close()
	jsonBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return 0, fmt.Errorf("error when converting response: %v", err)
	}
	var dollarExchangeResponse DollarExchangeResponse
	json.Unmarshal(jsonBytes, &dollarExchangeResponse)
	return dollarExchangeResponse.DollarPrice, nil
}
