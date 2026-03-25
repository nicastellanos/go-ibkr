package store

import (
	"sync"
)

// OptionData holds the Greeks and price for a specific contract
type OptionData struct {
	Price float64
	Delta float64
	Gamma float64
}

type MarketStore struct {
	sync.RWMutex
	// We use reqID as the key for fast lookups during callbacks
	Data map[int64]*OptionData
}

func NewMarketStore() *MarketStore {
	return &MarketStore{
		Data: make(map[int64]*OptionData),
	}
}

// UpdatePrice safely updates the last price for a contract
func (s *MarketStore) UpdatePrice(id int64, price float64) {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.Data[id]; !exists {
		s.Data[id] = &OptionData{}
	}
	s.Data[id].Price = price
}

// UpdateGreeks safely updates Delta and Gamma
func (s *MarketStore) UpdateGreeks(id int64, delta, gamma float64) {
	s.Lock()
	defer s.Unlock()
	if _, exists := s.Data[id]; !exists {
		s.Data[id] = &OptionData{}
	}
	s.Data[id].Delta = delta
	s.Data[id].Gamma = gamma
}

// GetSnapshot returns a copy of the current market state for the engine
func (s *MarketStore) GetSnapshot() map[int64]OptionData {
	s.RLock()
	defer s.RUnlock()
	
	snapshot := make(map[int64]OptionData)
	for id, val := range s.Data {
		snapshot[id] = *val
	}
	return snapshot
}