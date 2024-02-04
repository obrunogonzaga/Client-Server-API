package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type CurrencyQuote struct {
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

func main() {
	http.HandleFunc("/get-dollar-quote", GetDollarQuoteHandler)
	http.ListenAndServe(":8080", nil)
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
}

func GetDollarQuote(cep string) (*CurrencyQuote, error) {
	res, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var currencyQuote CurrencyQuote
	err = json.Unmarshal(body, &currencyQuote)
	return &currencyQuote, err
}
