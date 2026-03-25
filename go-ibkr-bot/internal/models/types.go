package models

// Position tracks how many units of an asset you actually hold
type Position struct {
	Symbol    string
	Quantity  float64
	IsOption  bool
	ReqID     int64 // Links to the MarketStore ID
}

// HedgeCommand tells the Order Manager what to do
type HedgeCommand struct {
	Symbol   string
	Quantity int // Negative for Sell, Positive for Buy
}