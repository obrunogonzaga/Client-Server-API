package main

import (
	"encoding/json"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
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
	ID        string `gorm:"primaryKey"`
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
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(dolarQuote)
	if err != nil {
		panic(err)
	}

	err = createCurrencyDB(dolarQuote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

func GetDollarQuote(cep string) (*CurrencyQuoteDTO, error) {
	res, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var currencyQuote CurrencyQuoteDTO
	err = json.Unmarshal(body, &currencyQuote)
	return &currencyQuote, err
}

func createCurrencyDB(dto *CurrencyQuoteDTO) error {

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
	db.Create(&CurrencyQuote{
		Value:     fBid,
		Code:      v.Code,
		Codein:    v.CodeIn,
		High:      fHigh,
		Low:       fLow,
		Name:      v.Name,
		PctChange: fPctChange,
	})

	return nil

}
