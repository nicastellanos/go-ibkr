package collector

import (
	"fmt"
	"go-ibkr-bot/internal/store"
	"github.com/scmhub/ibapi"
)

type IBKRCollector struct {
	ibapi.Wrapper
	Store *store.MarketStore
}

func NewIBKRCollector(s *store.MarketStore) *IBKRCollector {
	return &IBKRCollector{Store: s}
}

// TickPrice handles price updates (TickType 4 = Last Price)
func (w *IBKRCollector) TickPrice(reqId int64, tickType int64, price float64, attrib ibapi.TickAttrib) {
	fmt.Printf("DEBUG: Received Price %f for ID %d\n", price, reqId)
	if tickType == 4 && price > 0 {
		w.Store.UpdatePrice(reqId, price)
	}
}

func (w *IBKRCollector) TickOptionComputation(
	reqId ibapi.TickerID, 
	tickType ibapi.TickType, 
	tickAttrib int64,
	impliedVol float64, 
	delta float64, 
	optPrice float64, 
	pvDividend float64, 
	gamma float64, 
	vega float64, 
	theta float64, 
	undPrice float64,
) {
	// 13 = Model Computation (most stable delta)
	// We cast tickType to int64 to compare it if necessary, 
	// or just check against the ibapi constant.
	if int64(tickType) == 13 && delta > -2.0 {
		w.Store.UpdateGreeks(int64(reqId), delta, gamma)
	}
}

// Error handles system messages and connection issues
// Signature updated to match latest scmhub/ibapi requirements
func (w *IBKRCollector) Error(reqId int64, errorCode int64, advCode int64, errStr string, advMsg string) {
	if errorCode == 2104 || errorCode == 2106 || errorCode == 2158 {
		// These are just "Market Data Farm is connected" notifications
		return
	}
	fmt.Printf("IBKR Error [%d]: %s %s\n", errorCode, errStr, advMsg)
}

// NextValidId is called when the connection is fully established
func (w *IBKRCollector) NextValidId(reqId int64) {
	fmt.Printf("Connection established. Next valid Order ID: %d\n", reqId)
}