package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type CurrencyQuoteDTO struct {
	Dolar string `json:"dolar"`
}

func main() {
	quote, err := getDollarQuote()
	if err != nil {
		log.Printf("main() failed to get dollar quote: %s", err.Error())
		return
	}

	err = saveDollarQuoteinFile(*quote)
	if err != nil {
		log.Printf("main() failed to save dollar quote in file: %s", err.Error())
		return
	}
}

func getDollarQuote() (*CurrencyQuoteDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost:8080/get-dollar-quote", nil)
	if err != nil {
		fmt.Printf("main() failed to create request: %s", err.Error())
		return nil, fmt.Errorf("main() failed to create request: %s", err.Error())
	}

	res, err := http.DefaultClient.Do(req.WithContext(ctx))
	if err != nil {
		fmt.Printf("main() failed to do request: %s", err.Error())
		return nil, fmt.Errorf("main() failed to do request: %s", err.Error())
	}
	defer res.Body.Close()
	fmt.Println("Response status:", res.Status)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Printf("GetDollarQuote() failed to read response body: %s", err.Error())
		return nil, fmt.Errorf("GetDollarQuote() failed to read response body: %s", err.Error())
	}

	var currencyQuote CurrencyQuoteDTO
	err = json.Unmarshal(body, &currencyQuote)
	if err != nil {
		log.Printf("GetDollarQuote() failed to unmarshal response body: %s", err.Error())
		return nil, fmt.Errorf("GetDollarQuote() failed to unmarshal response body: %s", err.Error())
	}
	fmt.Println("Dollar quote:", currencyQuote.Dolar)
	return &currencyQuote, nil
}

func saveDollarQuoteinFile(quote CurrencyQuoteDTO) error {
	fileText := []byte("Dólar: " + quote.Dolar + "\n")

	err := os.WriteFile("cotacao.txt", fileText, 0644)
	if err != nil {
		return fmt.Errorf("saveDollarQuoteinFile() failed to write to file: %s", err.Error())
	}

	return nil
}
