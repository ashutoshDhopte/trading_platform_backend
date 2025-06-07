package routine

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"trading_platform_backend/db"
	"trading_platform_backend/orm"
)

type StockPriceGenerator struct {
	Ticker       string
	CurrentPrice int64
	MinPrice     int64
	MaxPrice     int64
	MaxChange    int64      // Maximum absolute change per minute
	mu           sync.Mutex // Mutex to protect CurrentPrice during concurrent access
}

func initStockPriceGenerator() {
	// Initialize the stock price generator
	go startGeneratorLoop()
}

func startGeneratorLoop() {

	stocks := db.GetAllStocks()
	stocksMap := make(map[string]*orm.Stocks)
	generators := make([]*StockPriceGenerator, 0)

	fmt.Printf("Initial Stock Price:")

	for i := range stocks {
		generators = append(generators, NewStockPriceGenerator(
			stocks[i].Ticker,
			stocks[i].OpeningPriceCents,
			stocks[i].MinPriceGeneratorCents,
			stocks[i].MaxPriceGeneratorCents,
			2.00,
		))
		fmt.Printf("%s $%d\n", stocks[i].Ticker, stocks[i].CurrentPriceCents)
		stocksMap[stocks[i].Ticker] = &stocks[i]
	}

	// Create a ticker that fires every minute
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop() // Ensure the ticker is stopped when main exits

	// Loop indefinitely, generating a new price every minute

	for range ticker.C {
		for _, generator := range generators {
			price := generator.GenerateNewPrice()
			//fmt.Printf("[%s] New Stock Price: %s $%d\n", time.Now().Format("15:04:05"), generator.Ticker, price)

			// Here you would typically publish this price, store it, or do something else with it.
			stocksMap[generator.Ticker].CurrentPriceCents = price
		}

		db.DB.Save(&stocks)

		WsHub.Broadcast <- ""

	}
}

func NewStockPriceGenerator(ticker string, openingPrice, minPrice, maxPrice, maxChange int64) *StockPriceGenerator {

	return &StockPriceGenerator{
		Ticker:       ticker,
		CurrentPrice: openingPrice,
		MinPrice:     minPrice,
		MaxPrice:     maxPrice,
		MaxChange:    maxChange,
	}
}

func (s *StockPriceGenerator) GenerateNewPrice() int64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	// IMPORTANT: Seed the random number generator only once in main, or use crypto/rand
	// For this example, we'll keep it here, but generally, seed once.
	// rand.Seed(time.Now().UnixNano()) // rand.Seed is deprecated in Go 1.20+, use rand.New(rand.NewSource(...))
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Generate a random change between -MaxChange and +MaxChange
	// We generate a float between -1.0 and 1.0, then scale it by MaxChange.
	// We then convert it to an int64. This conversion will naturally "round" or truncate.
	// To get a more random distribution of int64 changes, we can generate a random
	// int64 within the range directly.
	// Example: rand.Int63n(N) returns [0, N-1]
	// To get [-MaxChange, MaxChange]:
	// Generate random int in range [0, 2*MaxChange]
	// Subtract MaxChange to shift it to [-MaxChange, MaxChange]
	randomChange := r.Int63n(s.MaxChange*2+1) - s.MaxChange // This generates a random integer in [-MaxChange, MaxChange] inclusive.

	// Ensure minimum change magnitude if desired (e.g., at least 1 cent)
	// We only apply this if MaxChange is at least 1 cent (1).
	// If the generated change is 0 and we want at least 1 cent change,
	// we'll force it to 1 cent or -1 cent randomly.
	if s.MaxChange >= 1 && randomChange == 0 {
		if r.Float64() < 0.5 { // 50% chance to go up or down
			randomChange = 1
		} else {
			randomChange = -1
		}
	}

	newPrice := s.CurrentPrice + randomChange

	// Ensure the new price stays within the min/max range
	if newPrice < s.MinPrice {
		newPrice = s.MinPrice
	}
	if newPrice > s.MaxPrice {
		newPrice = s.MaxPrice
	}

	s.CurrentPrice = newPrice // No need for explicit rounding since it's int64

	return s.CurrentPrice
}
