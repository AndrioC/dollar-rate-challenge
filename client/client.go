package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type DollarRate struct {
	Bid string `json:"bid"`
}

type DollarData struct {
	USDBRL DollarRate `json:"USDBRL"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		println(err)
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	var dollarRate DollarData
	err = json.NewDecoder(res.Body).Decode(&dollarRate)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Dollar Bid: %s\n", dollarRate.USDBRL.Bid)
}
