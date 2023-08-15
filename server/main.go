package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type USDBRL struct {
	ID         int    `json:"USDBRLID"`
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type Cotacao struct {
	ID       int    `json:"CotacaoID"`
	USDBRLID int    `json:"USDBRLID" gorm:"foreignKey:USDBRLID;references:ID"`
	USDBRL   USDBRL `json:"USDBRL"`
	gorm.Model
}

func main() {
	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	ctxHttp, cancelHttp := context.WithTimeout(r.Context(), 200*time.Millisecond)
	defer cancelHttp()

	cotacaoRequest, err := http.NewRequestWithContext(ctxHttp, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
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

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Println(err)
	}
	db.AutoMigrate(&Cotacao{})

	var cotacao Cotacao
	err = json.Unmarshal(cotacaoResult, &cotacao)
	if err != nil {
		log.Println(err)
	}

	ctxDB, cancelDB := context.WithTimeout(r.Context(), 10*time.Millisecond)
	defer cancelDB()

	resultDB := db.WithContext(ctxDB).Create(&cotacao)

	if resultDB.Error != nil {
		log.Println(resultDB.Error)
	}

	w.Write([]byte(cotacao.USDBRL.Bid)) // O client.go precisará receber do server.go apenas o valor atual do câmbio (campo "bid" do JSON)
}
