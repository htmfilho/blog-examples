package main

import "testing"

var stocks = []*Stock{
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

func BenchmarkGetPricesInSequence(b *testing.B) {
	for n := 0; n < b.N; n++ {
		stocks = getPricesInSequence(stocks)
	}
}

func BenchmarkGetPricesInParallel(b *testing.B) {
	for n := 0; n < b.N; n++ {
		stocks = getPricesInParallel(stocks)
	}
}
