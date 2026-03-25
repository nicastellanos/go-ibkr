package engine

import (
	"fmt"
	"go-ibkr-bot/internal/models"
	"go-ibkr-bot/internal/store"
	"time"
)

type Rebalancer struct {
	Store           *store.MarketStore
	Positions       []models.Position
	HedgeThreshold  float64 // Rebalance if Delta > 10.0, for example
	CheckInterval   time.Duration
}

func NewRebalancer(s *store.MarketStore, threshold float64) *Rebalancer {
	return &Rebalancer{
		Store:          s,
		HedgeThreshold: threshold,
		CheckInterval:  1 * time.Second,
	}
}

func (r *Rebalancer) Start() {
	ticker := time.NewTicker(r.CheckInterval)
	for range ticker.C {
		r.calculateAndHedge()
	}
}

func (r *Rebalancer) calculateAndHedge() {
	var totalDelta float64
	snapshot := r.Store.GetSnapshot()

	for _, pos := range r.Positions {
		data, exists := snapshot[pos.ReqID]
		if !exists {
			continue
		}

		if pos.IsOption {
			// Options have a multiplier (usually 100)
			totalDelta += (data.Delta * pos.Quantity * 100)
		} else {
			// Stock delta is always 1.0 per share
			totalDelta += (1.0 * pos.Quantity)
		}
	}

	fmt.Printf("Current Portfolio Delta: %.2f\n", totalDelta)

	// Decision Logic
	if totalDelta > r.HedgeThreshold {
		fmt.Printf("ACTION: Sell %d shares to hedge\n", int(totalDelta))
	} else if totalDelta < -r.HedgeThreshold {
		fmt.Printf("ACTION: Buy %d shares to hedge\n", int(-totalDelta))
	}
}