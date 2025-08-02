package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Fatalf("Erro criando requisição: %v", err)
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Timeout ao receber resposta do servidor")
		}
		log.Fatalf("Erro fazendo requisição: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Fatalf("Status inesperado do servidor: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Erro lendo resposta: %v", err)
	}

	var data map[string]string
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatalf("Erro ao decodificar JSON: %v", err)
	}

	bid, ok := data["bid"]
	if !ok {
		log.Fatalf("Campo bid não encontrado na resposta")
	}

	err = os.WriteFile("cotacao.txt", []byte(fmt.Sprintf("Dólar: %s", bid)), 0644)
	if err != nil {
		log.Fatalf("Erro ao salvar arquivo: %v", err)
	}

	fmt.Printf("Cotação salva em cotacao.txt: Dólar: %s\n", bid)
}
