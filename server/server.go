package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
}

func DollarRateQuery() (*DollarData, error) {
	req, err := http.Get(API_URL)

	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	res, err := io.ReadAll(req.Body)

	if err != nil {
		panic(err)
	}

	var data DollarData
	fmt.Println(string(res))

	err = json.Unmarshal(res, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil

}
