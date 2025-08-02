package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type CotacaoAPIResponse struct {
	USDBRL struct {
		Bid string `json:"bid"`
	} `json:"USDBRL"`
}

func main() {
	db, err := openDbConnection()
	if err != nil {
		log.Fatalf("Erro ao abrir banco: %v", err)
	}
	defer db.Close()

	http.HandleFunc("/cotacao", cotacaoHandler(db))

	log.Println("Server rodando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func openDbConnection() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./data/cotacoes.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS cotacoes (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		bid TEXT,
		timestamp DATETIME DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func cotacaoHandler(db *sql.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctxApi, cancelApi := context.WithTimeout(r.Context(), 200*time.Millisecond)
		defer cancelApi()

		bid, err := fetchCotacao(ctxApi)
		if err != nil {
			http.Error(w, "Erro ao obter cotação: "+err.Error(), http.StatusInternalServerError)
			log.Printf("Erro fetchCotacao: %v", err)
			return
		}

		ctxDB, cancelDB := context.WithTimeout(r.Context(), 10*time.Millisecond)
		defer cancelDB()

		err = salvarCotacao(ctxDB, db, bid)
		if err != nil {
			log.Printf("Erro salvarCotacao: %v", err)
		}

		resp := map[string]string{"bid": bid}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

func fetchCotacao(ctx context.Context) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return "", err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("Timeout ao chamar API de cotação")
		}
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status da API: %d", resp.StatusCode)
	}

	var data CotacaoAPIResponse
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return "", err
	}

	return data.USDBRL.Bid, nil
}

func salvarCotacao(ctx context.Context, db *sql.DB, bid string) error {
	done := make(chan error, 1)
	go func() {
		_, err := db.ExecContext(ctx, "INSERT INTO cotacoes(bid) VALUES(?)", bid)
		done <- err
	}()

	select {
	case <-ctx.Done():
		log.Println("Timeout ao salvar cotação no banco")
		return ctx.Err()
	case err := <-done:
		return err
	}
}
