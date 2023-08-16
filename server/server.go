package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

	dollarData, err := DollarRateQuery(ctx)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	dollarDataJSON, err := json.Marshal(dollarData)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = InsertDollarRateRequest(dollarData.USDBRL.Bid)

	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(dollarDataJSON)

}

func DollarRateQuery(ctx context.Context) (*DollarData, error) {
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)

	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "GET", API_URL, nil)

	if err != nil {
		fmt.Printf("Error while creating request: %v", err)
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		fmt.Printf("Error while sending request: %v", err)
		return nil, err
	}

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

func InsertDollarRateRequest(bid string) error {
	dbPath := "./db/main.db"

	_, err := os.Stat(dbPath)

	if os.IsNotExist(err) {
		fmt.Println("database file not found, creating...")

		err = os.MkdirAll("./db", os.ModePerm)

		if err != nil {
			return err
		}

		file, err := os.Create(dbPath)

		if err != nil {
			return err
		}

		file.Close()
	}

	insertDataCtx, insertDataCancel := context.WithTimeout(context.Background(), 10*time.Millisecond)

	defer insertDataCancel()

	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		fmt.Printf("Error trying to create database: %v\n", err)
		return err
	}

	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS dollarraterequests (
						id INTEGER PRIMARY KEY AUTOINCREMENT,
						bid TEXT,
						created_at DATE DEFAULT CURRENT_TIMESTAMP
					)`)

	if err != nil {
		fmt.Printf("Error while creating table: %v\n", err)
		return err
	}

	stmt, err := db.PrepareContext(insertDataCtx, "INSERT INTO dollarraterequests (bid) VALUES (?)")

	if err != nil {
		fmt.Printf("Error inserting data: %v\n", err)
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(bid)
	if err != nil {
		fmt.Printf("Error while trying to save data: %v", err)
		return err
	}

	return nil

}
