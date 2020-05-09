package models

type CoinGeckoCoinListArray []struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

type CoinGeckoMarketData struct {
	CurrentPrice map[string]float64 `json:"current_price"`
}
type CoinGeckoCoinInfo struct {
	ID         string              `json:"id"`
	Symbol     string              `json:"symbol"`
	Name       string              `json:"name"`
	MarketData CoinGeckoMarketData `json:"market_data"`
}
