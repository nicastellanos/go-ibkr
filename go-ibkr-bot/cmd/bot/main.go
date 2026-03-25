package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go-ibkr-bot/internal/collector"
	"go-ibkr-bot/internal/engine"
	"go-ibkr-bot/internal/models"
	"go-ibkr-bot/internal/store"

	"github.com/scmhub/ibapi"
)

func main() {
	// 1. Initialize the Shared State
	marketStore := store.NewMarketStore()

	portfolio := []models.Position{
    {Symbol: "ES", Quantity: 1, IsOption: false, ReqID: 10},  // E-mini S&P 500
    {Symbol: "RTY", Quantity: 1, IsOption: false, ReqID: 11}, // E-mini Russell 2000
}

	// 3. Initialize the Collector and Rebalancer
	ibkrWrapper := collector.NewIBKRCollector(marketStore)
	rebalancer := engine.NewRebalancer(marketStore, 10.0) // Hedge if delta > 10
	rebalancer.Positions = portfolio

	// 4. Setup and Connect the IBKR Client
	client := ibapi.NewEClient(ibkrWrapper)
	err := client.Connect("127.0.0.1", 7497, 999)
	if err != nil {
		log.Fatalf("Could not connect to TWS: %v", err)
	}


	for _, pos := range portfolio {
		contract := &ibapi.Contract{
			Symbol:   pos.Symbol,
			Currency: "USD",
		}

		switch pos.Symbol {
			case "ES", "RTY":
				contract.SecType = "FUT"
				contract.Exchange = "CME" // Or "GLOBEX"
				contract.LastTradeDateOrContractMonth = "202606" // June 2026 Expiry
				contract.Multiplier = "50"
				
			case "N225":
				contract.SecType = "IND"
				contract.Exchange = "OSE"
				contract.Currency = "JPY"

			default:
				contract.SecType = "STK"
				contract.Exchange = "SMART"
		}

		client.ReqMktData(pos.ReqID, contract, "", false, false, nil)
	}

	// 6. Start the Rebalancer in its own Goroutine
	go rebalancer.Start()

	// 7. Graceful Shutdown Handling
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	fmt.Println("Bot is live. Monitoring Delta Neutrality...")
	<-stop

	fmt.Println("\nShutting down safely...")
	client.Disconnect()
}