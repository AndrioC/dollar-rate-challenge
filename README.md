# Client-Server-API Challenge

In this challenge, concepts of HTTP web server, contexts, database, and file manipulation in Go will be applied.

There are two programs in this repository:

**client.go**

**server.go**

**Challenge requirements:**

1. The client.go program should make an HTTP request to server.go, requesting the current exchange rate for the US Dollar.

2. The server.go program should consume the API containing the USD to BRL exchange rate at the following address: https://economia.awesomeapi.com.br/json/last/USD-BRL, and then return the result in JSON format to the client.

3. Using the "context" package, server.go should record each received exchange rate in the SQLite database. The maximum timeout for calling the USD exchange rate API should be 200ms, and the maximum timeout for persisting the data in the database should be 10ms.

4. The client.go program should only receive the current exchange rate value ("bid" field from the JSON) from server.go. Using the "context" package, client.go will have a maximum timeout of 300ms to receive the result from server.go.

5. All 3 contexts should log an error if the execution time is insufficient.

6. The client.go program should save the current exchange rate in a file named "cotacao.txt" in the format: Dólar: {value}

7. The required endpoint generated by server.go for this challenge will be: "/cotacao", and the HTTP server should use port 8080.

## How to run this project?

First of all, you need to install the dependency packages, run:

**go mod tidy**

Now you need to run both programs in separate terminals and in the following order:

1. Run go run server/server.go
2. Run go run client/client.go

The dollar rate value will be saved in db/main.db, and a file named cotacao.txt will be created.
