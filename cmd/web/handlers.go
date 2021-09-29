package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/asbeeq/binance/pkg/client_http"
)

func doRequest(done <-chan bool, data chan<- *client_http.OrderBook, ticker *time.Ticker, symbol string) {
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			fmt.Println("Tick at", t)

			orderBooks, err := client_http.MakeRequest(symbol)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// положить в канал orderBooks
			data <- orderBooks

			// отображение на консоли
			fmt.Printf("Bids: %v\n\nAsks: %v\n\nSum Bids Quantity: %f\n\nSum Asks Quantity: %f\n",
				orderBooks.Bids, orderBooks.Asks, orderBooks.SumBidsQuantity, orderBooks.SumAsksQuantity)
		}
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// files := []string{
	// 	"./ui/html/home.page.tmpl",
	// 	"./ui/html/base.layout.tmpl",
	// }

	// ts, err := template.ParseFiles(files...)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	http.Error(w, "Internal Server Error", 500)
	// 	return
	// }

	// err = ts.Execute(w, nil)
	// if err != nil {
	// 	log.Println(err.Error())
	// 	http.Error(w, "Internal Server Error", 500)
	// }

	symbol := "ETHBTC" // r.URL.Query().Get("symbol")
	if symbol != "" {
		ticker := time.NewTicker(5 * time.Second)
		done := make(chan bool)
		data := make(chan *client_http.OrderBook)

		go doRequest(done, data, ticker, symbol)

	outer:
		for {
			select {
			case <-r.Context().Done():
				close(done)
				break outer
			case <-data:
			}
		}
	}

}
