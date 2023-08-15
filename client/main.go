package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8081", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	ctxHttp, cancelHttp := context.WithTimeout(r.Context(), 300*time.Millisecond)
	defer cancelHttp()

	cotacaoRequest, err := http.NewRequestWithContext(ctxHttp, "GET", "http://localhost:8080/cotacao", nil)
	if err != nil {
		log.Println(err)
	}

	cotacaoResponse, err := http.DefaultClient.Do(cotacaoRequest)
	if err != nil {
		log.Println(err)
	}
	defer cotacaoResponse.Body.Close()

	cotacaoResult, err := io.ReadAll(cotacaoResponse.Body)
	if err != nil {
		log.Println(err)
	}

	if fileNotExists() {
		createFile(cotacaoResult)
	} else {
		editFile(cotacaoResult)
	}

	w.Write([]byte(cotacaoResult))
}

func fileNotExists() bool {
	if _, err := os.Stat("cotacao.txt"); err != nil {
		return true
	}
	return false
}

func createFile(cotacaoResult []byte) {
	fileCotacao, err := os.Create("cotacao.txt")
	defer fileCotacao.Close()

	if err != nil {
		log.Println(err)
	}

	tamanho, err := fileCotacao.WriteString("Dólar: " + string(cotacaoResult) + "\n")
	if err != nil {
		log.Println(err)
		log.Println(tamanho)
	}
}

func editFile(cotacaoResult []byte) {
	fileCotacao, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_WRONLY, 0644)
	defer fileCotacao.Close()

	if err != nil {
		log.Println(err)
	}

	tamanho, err := fileCotacao.WriteString("Dólar: " + string(cotacaoResult) + "\n")
	if err != nil {
		log.Println(err)
		log.Println(tamanho)
	}
}
