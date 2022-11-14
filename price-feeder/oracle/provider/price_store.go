package provider

import (
	"fmt"
	"sync"

	"github.com/umee-network/umee/price-feeder/oracle/types"
)

// PriceStore is an embedded struct in each provider that manages the in memory
// store of subscribed currency pairs, candles prices, and ticker prices. It also
// handles thread safety and pruning of old candle prices.
type PriceStore struct {
	tickers         map[string]types.TickerPrice
	candles         map[string][]types.CandlePrice
	subscribedPairs map[string]types.CurrencyPair

	subscribedPairsMtx sync.RWMutex
	tickerMtx          sync.RWMutex
	candleMtx          sync.RWMutex
}

func NewPriceStore() PriceStore {
	return PriceStore{
		tickers:         map[string]types.TickerPrice{},
		candles:         map[string][]types.CandlePrice{},
		subscribedPairs: map[string]types.CurrencyPair{},
	}
}

// GetTickerPrices returns the tickerPrices based on the provided pairs. Returns an
// error if ANY of the currency pairs are not available.
func (ps *PriceStore) GetTickerPrices(pairs ...types.CurrencyPair) (map[string]types.TickerPrice, error) {
	ps.tickerMtx.RLock()
	defer ps.tickerMtx.RUnlock()

	tickerPrices := make(map[string]types.TickerPrice, len(pairs))
	for _, cp := range pairs {
		key := cp.String()
		ticker, ok := ps.tickers[key]
		if !ok {
			return nil, fmt.Errorf("failed to get ticker price for %s", key)
		}
		tickerPrices[key] = ticker
	}
	return tickerPrices, nil
}

// GetCandlePrices returns a copy of the the candlePrices based on the provided pairs.
// Returns an error if ANY of the currency pairs are not available
func (ps *PriceStore) GetCandlePrices(pairs ...types.CurrencyPair) (map[string][]types.CandlePrice, error) {
	ps.candleMtx.RLock()
	defer ps.candleMtx.RUnlock()

	candlePrices := make(map[string][]types.CandlePrice, len(pairs))
	for _, cp := range pairs {
		key := cp.String()
		candles, ok := ps.candles[key]
		if !ok {
			return nil, fmt.Errorf("failed to get candle prices for %s", key)
		}
		candlesCopy := make([]types.CandlePrice, 0, len(candles))
		candlesCopy = append(candlesCopy, candles...)
		candlePrices[key] = candlesCopy
	}
	return candlePrices, nil
}

func (ps *PriceStore) setTickerPair(ticker types.TickerPrice, currencyPair string) {
	ps.tickerMtx.Lock()
	defer ps.tickerMtx.Unlock()

	ps.tickers[currencyPair] = ticker
}

func (ps *PriceStore) setCandlePair(candle types.CandlePrice, currencyPair string) {
	ps.candleMtx.Lock()
	defer ps.candleMtx.Unlock()

	staleTime := PastUnixTime(providerCandlePeriod)
	newCandles := []types.CandlePrice{candle}

	for _, c := range ps.candles[currencyPair] {
		if staleTime < c.TimeStamp {
			newCandles = append(newCandles, c)
		}
	}
	ps.candles[currencyPair] = newCandles
}

// setSubscribedPairs sets N currency pairs to the map of subscribed pairs.
func (ps *PriceStore) setSubscribedPairs(cps ...types.CurrencyPair) {
	ps.subscribedPairsMtx.Lock()
	defer ps.subscribedPairsMtx.Unlock()

	for _, cp := range cps {
		ps.subscribedPairs[cp.String()] = cp
	}
}

// AddSubscribedPairs adds any unique currency pairs to the subscribed currency
// pairs map and returns the pairs added with the duplicates removed.
func (ps *PriceStore) addSubscribedPairs(cps ...types.CurrencyPair) []types.CurrencyPair {
	newPairs := []types.CurrencyPair{}
	for _, cp := range cps {
		if ps.isSubscribed(cp.String()) {
			newPairs = append(newPairs, cp)
		}
	}
	ps.setSubscribedPairs(newPairs...)
	return newPairs
}

func (ps *PriceStore) isSubscribed(currencyPair string) bool {
	ps.subscribedPairsMtx.RLock()
	defer ps.subscribedPairsMtx.RUnlock()

	if _, ok := ps.subscribedPairs[currencyPair]; ok {
		return true
	}
	return false
}
