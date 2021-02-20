package offer

import (
	"encoding/json"
	"net/http"
	"strconv"
)

// GetOffer GET: /offer?savings=&listingPrice=&downPayment=&closingCosts=
func GetOffer(res http.ResponseWriter, req *http.Request) {
	type Offer struct {
		Savings      float32 `json:"savings"`
		ListingPrice float32 `json:"listingPrice"`
		DownPayment  float32 `json:"downPayment"`
		ClosingCosts float32 `json:"closingCosts"`
		MaximumBid   float32 `json:"maximumBid"`
		Margin       float32 `json:"margin"`
	}

	savings := getURLParam(req, "savings")
	listingPrice := getURLParam(req, "listingPrice")
	downPayment := getURLParam(req, "downPayment")
	closingCosts := getURLParam(req, "closingCosts")

	maxBid := CalcMaxBid(savings, listingPrice, downPayment, closingCosts)
	margin := CalcMargin(savings, listingPrice, downPayment, closingCosts)

	offer := &Offer{
		Savings:      savings,
		ListingPrice: listingPrice,
		DownPayment:  downPayment,
		ClosingCosts: closingCosts,
		MaximumBid:   maxBid,
		Margin:       margin,
	}

	res.Header().Set("Content-Type", "application/json; charset=UTF-8")
	err := json.NewEncoder(res).Encode(offer)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func getURLParam(req *http.Request, param string) float32 {
	value, err := strconv.ParseFloat(req.URL.Query()[param][0], 32)
	if err != nil {
		return 0.0
	}
	return float32(value)
}
