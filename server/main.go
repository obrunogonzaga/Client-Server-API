package main

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type CurrencyQuoteDTO struct {
	USDBRL struct {
		Ask        string `json:"ask"`
		Bid        string `json:"bid"`
		Code       string `json:"code"`
		CodeIn     string `json:"codein"`
		CreateDate string `json:"create_date"`
		High       string `json:"high"`
		Low        string `json:"low"`
		Name       string `json:"name"`
		PctChange  string `json:"pctChange"`
		TimesStamp string `json:"timestamp"`
	} `json:"USDBRL"`
}

type CurrencyQuote struct {
	gorm.Model
	ID        string `gorm:"type:uuid;"`
	Value     float64
	Code      string
	Codein    string
	High      float64
	Low       float64
	Name      string
	PctChange float64
}

func main() {

	http.HandleFunc("/get-dollar-quote", GetDollarQuoteHandler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}

}

func GetDollarQuoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/get-dollar-quote" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	dolarQuote, err := GetDollarQuote("USD-BRL")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error: %s", err.Error())
		return
	}

	err = createCurrencyDB(dolarQuote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error: %s", err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dolarQuote)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return
	}
}

func GetDollarQuote(cep string) (*CurrencyQuoteDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, err
	}

	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return nil, err
	}
	var currencyQuote CurrencyQuoteDTO
	err = json.Unmarshal(body, &currencyQuote)
	return &currencyQuote, err
}

func createCurrencyDB(dto *CurrencyQuoteDTO) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		err := os.Mkdir("./data", os.ModePerm)
		if err != nil {
			return err
		}
	}

	db, err := gorm.Open(sqlite.Open("./data/currency.db"), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&CurrencyQuote{})
	if err != nil {
		return err
	}

	v := dto.USDBRL
	fBid, err := strconv.ParseFloat(v.Bid, 64)
	if err != nil {
		return err
	}

	fHigh, err := strconv.ParseFloat(v.High, 64)
	if err != nil {
		return err
	}

	fLow, err := strconv.ParseFloat(v.Low, 64)
	if err != nil {
		return err
	}

	fPctChange, err := strconv.ParseFloat(v.PctChange, 64)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}

	err = db.WithContext(ctx).Create(&CurrencyQuote{
		ID:        uuid.New().String(),
		Value:     fBid,
		Code:      v.Code,
		Codein:    v.CodeIn,
		High:      fHigh,
		Low:       fLow,
		Name:      v.Name,
		PctChange: fPctChange,
	}).Error
	if err != nil {
		log.Printf("Error: %s", err.Error())
		return err
	}

	return nil
}
