package models

// BestOrderBook struct
type BestOrderBook struct {
	Ask Order `json:"ask"` //asks.Price > any bids.Price
	Bid Order `json:"bid"`
}

// Order struct
type Order struct {
	Amount float64 `json:"amount"`
	Price  float64 `json:"price"`
}
