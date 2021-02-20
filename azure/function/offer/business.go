package offer

func CalcMaxBid(savings, listingPrice, downPayment, closingCosts float32) float32 {
	return listingPrice + CalcMargin(savings, listingPrice, downPayment, closingCosts)
}

func CalcMargin(savings, listingPrice, downPayment, closingCosts float32) float32 {
	return savings - (listingPrice * (downPayment / 100.0)) - closingCosts
}