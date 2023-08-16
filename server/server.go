package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

type DollarRate struct {
	Bid string `json:"bid"`
}

type DollarData struct {
	USDBRL DollarRate `json:"USDBRL"`
}

const API_URL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

func main() {
	http.HandleFunc("/cotacao", DollarRateQueryHandler)
	http.ListenAndServe(":8080", nil)
}

func DollarRateQueryHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log.Println("Making request..")

	defer log.Println("Request finished")
	if r.URL.Path == "/" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	dollarRate, err := DollarRateQuery(ctx)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dollarRate)

}

func DollarRateQuery(ctx context.Context) (*DollarData, error) {
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)

	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", API_URL, nil)

	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return nil, err
	}

	var data DollarData

	err = json.Unmarshal(body, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil

}
