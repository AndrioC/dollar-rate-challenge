package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type DollarRate struct {
	Bid string `json:"bid"`
}

type DollarData struct {
	USDBRL DollarRate `json:"USDBRL"`
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		fmt.Printf("Error while creating request: %v", err)
		panic(err)
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("Error while sending request: %v", err)
		panic(err)
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if len(body) == 0 {
		fmt.Println("The response from server is empty")
		return
	}

	if err != nil {
		panic(err)
	}

	var dollarData DollarData

	err = json.Unmarshal(body, &dollarData)
	if err != nil {
		panic(err)
	}

	err = SaveDollarRateIntoFile(dollarData.USDBRL.Bid)

	if err != nil {
		panic(err)
	}
}

func SaveDollarRateIntoFile(bid string) error {
	f, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}

	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("Dólar: %v\n", bid))

	if err != nil {
		fmt.Fprintf(os.Stderr, "Erro ao escrever no arquivo: %v\n", err)
		return err
	}

	return nil
}
