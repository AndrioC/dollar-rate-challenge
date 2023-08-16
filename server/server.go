package main

import (
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

	dollarRate, err := DollarRateQuery()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dollarRate)

	select {
	case <-time.After(3 * time.Second):
		log.Printf("Request successfully completed")
		w.Write([]byte("Request successfully completed"))

	case <-ctx.Done():
		log.Printf("Request canceled by client")
		w.Write([]byte("Request canceled by client"))
	}
}

func DollarRateQuery() (*DollarData, error) {
	req, err := http.Get(API_URL)

	if err != nil {
		return nil, err
	}

	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)

	if err != nil {
		return nil, err
	}

	var data DollarData

	err = json.Unmarshal(res, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil

}
