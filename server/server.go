package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type DollarRate struct {
	Bid string `json:"bid"`
}

type DollarData struct {
	USDBRL DollarRate `json:"USDBRL"`
}

const API_URL = "https://economia.awesomeapi.com.br/json/last/USD-BRL"

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
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
		fmt.Fprintf(os.Stderr, "Erro ao fazer parse da resposta: %v\n", err)
	}

	fmt.Println(data.USDBRL.Bid)
}
