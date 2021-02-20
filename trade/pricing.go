package main

import (
	"fmt"
	"sync"
	"time"
)

// Supported exchanges
const (
	NASDAQ = "NASDAQ"
	NYSE   = "NYSE"
	TSX    = "TSX"
)

// Stock represents an investment
type Stock struct {
	Exchange string
	Ticker   string
	Name     string
	Price    int32
}

func main() {
	stocks := []*Stock{
		{
			Exchange: NASDAQ,
			Ticker:   "GOOGL",
			Name:     "Alphabet Inc.",
		}, {
			Exchange: NYSE,
			Ticker:   "GE",
			Name:     "General Electric CO",
		}, {
			Exchange: TSX,
			Ticker:   "SHOP",
			Name:     "Shopify Inc.",
		}, {
			Exchange: NASDAQ,
			Ticker:   "AAPL",
			Name:     "Apple Inc.",
		}, {
			Exchange: TSX,
			Ticker:   "BMO",
			Name:     "Bank of Montreal",
		}, {
			Exchange: TSX,
			Ticker:   "CP",
			Name:     "Canadian Pacific Raiway Ltd.",
		},
	}

	stocks = getPricesInSequence(stocks)

	for _, stock := range stocks {
		fmt.Printf("%s: %d\n", stock.Ticker, stock.Price)
	}
}

func getPricesInSequence(stocks []*Stock) []*Stock {
	exchangeFactory := prepareExchangeStrategies(stocks)
	return searchStocksInSequence(exchangeFactory, stocks)
}

func getPricesInParallel(stocks []*Stock) []*Stock {
	exchangeFactory := prepareExchangeStrategies(stocks)
	return searchStocksInParallel(exchangeFactory, stocks)
}

func prepareExchangeStrategies(stocks []*Stock) *ExchangeFactory {
	exchangeFactory := &ExchangeFactory{}
	for _, stock := range stocks {
		pricing := exchangeFactory.GetExchangePricing(stock)
		pricing.addStock(stock)
	}
	return exchangeFactory
}

func searchStocksInSequence(exchangeFactory *ExchangeFactory, stocks []*Stock) []*Stock {
	for _, exchange := range exchangeFactory.exchanges {
		err := exchange.search()
		if err != nil {
			fmt.Printf("Couldn't fetch prices from %s.\n", exchange.getName())
		}
	}
	return stocks
}

func searchStocksInParallel(exchangeFactory *ExchangeFactory, stocks []*Stock) []*Stock {
	errors := make(chan error, 2)
	var wg sync.WaitGroup
	for _, exchange := range exchangeFactory.exchanges {
		wg.Add(1)
		go func(exchange Pricing) {
			defer wg.Done()
			err := exchange.search()
			if err != nil {
				errors <- err
				fmt.Printf("Couldn't fetch prices from %s.\n", exchange.getName())
				return
			}
		}(exchange)
	}
	wg.Wait()
	close(errors)

	return stocks
}

// Pricing shapes the strategy pattern to be applied to specialized structs.
type Pricing interface {
	addStock(stock *Stock)
	search() error
	getName() string
}

// ExchangeFactory holds the pricing strategies for the duration of the execution.
type ExchangeFactory struct {
	exchanges map[string]Pricing
}

// GetExchangePricing returns the Pricing based on the stock exchange.
func (ec *ExchangeFactory) GetExchangePricing(stock *Stock) Pricing {
	if ec.exchanges == nil {
		ec.exchanges = make(map[string]Pricing)
		ec.exchanges[NASDAQ] = new(NasdaqPricing)
		ec.exchanges[NYSE] = new(NYSEPricing)
		ec.exchanges[TSX] = new(TsxPricing)
	}
	return ec.exchanges[stock.Exchange]
}

// NasdaqPricing is the strategy implementation to get prices from Nasdaq.
type NasdaqPricing struct {
	stocks []*Stock
}

func (np *NasdaqPricing) addStock(stock *Stock) {
	np.stocks = append(np.stocks, stock)
}

func (np *NasdaqPricing) search() error {
	if len(np.stocks) > 0 {
		time.Sleep(time.Millisecond)
		np.stocks[0].Price = 2344500
		np.stocks[1].Price = 5439990
	}

	return nil
}

func (np *NasdaqPricing) getName() string {
	return "Nasdaq"
}

// NYSEPricing is the strategy implementation to get prices from NYSE.
type NYSEPricing struct {
	stocks []*Stock
}

func (nyp *NYSEPricing) addStock(stock *Stock) {
	nyp.stocks = append(nyp.stocks, stock)
}

func (nyp *NYSEPricing) search() error {
	if len(nyp.stocks) > 0 {
		time.Sleep(time.Millisecond)
		nyp.stocks[0].Price = 344500
	}

	return nil
}

func (nyp *NYSEPricing) getName() string {
	return "New York Stock Exchange"
}

// TsxPricing is the strategy implementation to get prices from Tsx.
type TsxPricing struct {
	stocks []*Stock
}

func (tp *TsxPricing) addStock(stock *Stock) {
	tp.stocks = append(tp.stocks, stock)
}

func (tp *TsxPricing) search() error {
	if len(tp.stocks) > 0 {
		time.Sleep(time.Millisecond)
		tp.stocks[0].Price = 8344500
		tp.stocks[1].Price = 239990
		tp.stocks[2].Price = 39990
	}

	return nil
}

func (tp *TsxPricing) getName() string {
	return "Toronto Stock Exchange"
}
