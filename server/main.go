package main

import (
	"context"
	"encoding/json"
	"fmt"
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

type CurrencyQuoteDAO struct {
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

type CurrencyQuoteDTO struct {
	Dolar string `json:"dolar"`
}

type App struct {
	DB *gorm.DB
}

func main() {

	app, err := InitializeApp()
	if err != nil {
		fmt.Printf("main() failed to connect database: %s", err.Error())
		return
	}

	http.HandleFunc("/cotacao", app.GetDollarQuoteHandler)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Printf("main() failed to start server: %s", err.Error())
		return
	}

}

func (app *App) GetDollarQuoteHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cotacao" {
		http.Error(w, "404 not found.", http.StatusNotFound)
		return
	}

	dolarQuote, err := GetDollarQuote("USD-BRL")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("GetDollarQuoteHandler() failed to get dollar quote: %s", err.Error())
		return
	}

	err = app.createCurrencyDB(dolarQuote)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("GetDollarQuoteHandler() failed to create currency db: %s", err.Error())
		return
	}

	mountResponse(w, dolarQuote)

}

func mountResponse(w http.ResponseWriter, quote *CurrencyQuoteDAO) {
	dolarQuote := &CurrencyQuoteDTO{
		Dolar: quote.USDBRL.Bid,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err := json.NewEncoder(w).Encode(dolarQuote)
	if err != nil {
		log.Printf("mountResponse() failed to encode response: %s", err.Error())
		return
	}
}

func GetDollarQuote(cep string) (*CurrencyQuoteDAO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3000*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		log.Printf("GetDollarQuote() failed to create request: %s", err.Error())
		return nil, err
	}

	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		log.Printf("GetDollarQuote() failed to do request: %s", err.Error())
		return nil, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("GetDollarQuote() failed to read response body: %s", err.Error())
		return nil, err
	}
	var currencyQuote CurrencyQuoteDAO
	err = json.Unmarshal(body, &currencyQuote)
	return &currencyQuote, err
}

func InitializeApp() (*App, error) {
	if _, err := os.Stat("./data"); os.IsNotExist(err) {
		err := os.Mkdir("./data", os.ModePerm)
		if err != nil {
			log.Printf("connectDatadase() failed to create data folder: %s", err.Error())
			return nil, fmt.Errorf("connectDatadase() failed to create data folder: %s", err.Error())
		}
	}

	db, err := gorm.Open(sqlite.Open("./data/currency.db"), &gorm.Config{})
	if err != nil {
		log.Printf("connectDatadase() failed to open db: %s", err.Error())
		return nil, fmt.Errorf("connectDatadase() failed to open db: %s", err.Error())
	}

	err = db.AutoMigrate(&CurrencyQuote{})
	if err != nil {
		log.Printf("connectDatadase() failed to migrate db: %s", err.Error())
		return nil, fmt.Errorf("connectDatadase() failed to migrate db: %s", err.Error())
	}

	return &App{DB: db}, nil
}

func (app *App) createCurrencyDB(dto *CurrencyQuoteDAO) error {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	v := dto.USDBRL
	fBid, err := strconv.ParseFloat(v.Bid, 64)
	if err != nil {
		log.Printf("createCurrencyDB() failed to parse float: %s", err.Error())
		return fmt.Errorf("createCurrencyDB() failed to parse float: %s", err.Error())
	}

	fHigh, err := strconv.ParseFloat(v.High, 64)
	if err != nil {
		log.Printf("createCurrencyDB() failed to parse float: %s", err.Error())
		return fmt.Errorf("createCurrencyDB() failed to parse float: %s", err.Error())
	}

	fLow, err := strconv.ParseFloat(v.Low, 64)
	if err != nil {
		log.Printf("createCurrencyDB() failed to parse float: %s", err.Error())
		return fmt.Errorf("createCurrencyDB() failed to parse float: %s", err.Error())
	}

	fPctChange, err := strconv.ParseFloat(v.PctChange, 64)
	if err != nil {
		log.Printf("createCurrencyDB() failed to parse float: %s", err.Error())
		return fmt.Errorf("createCurrencyDB() failed to parse float: %s", err.Error())
	}

	err = app.DB.WithContext(ctx).Create(&CurrencyQuote{
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
		log.Printf("createCurrencyDB() failed to create currency quote: %s", err.Error())
		return fmt.Errorf("createCurrencyDB() failed to create currency quote: %s", err.Error())
	}

	return nil
}
