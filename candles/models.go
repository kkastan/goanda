package candles

// CandleResponse represents the top level Oanda CandlestickResponse model.
// See https://developer.oanda.com/rest-live-v20/instrument-df/
type CandleResponse struct {
	Instrument  string   `json:"instrument"`
	Granularity string   `json:"granularity"`
	Candles     []Candle `json:"candles"`
}

// Candle represents the Oanda Candlestick model.
// See https://developer.oanda.com/rest-live-v20/instrument-df/
type Candle struct {
	Complete bool             `json:"complete"`
	Volume   int64            `json:"volume"`
	Time     string           `json:"time"`
	Bid      *CandleStickData `json:"bid"`
	Ask      *CandleStickData `json:"ask"`
	Mid      *CandleStickData `json:"mid"`
}

// CandleStickData represents the Oanda CandlestickData model.
// See https://developer.oanda.com/rest-live-v20/instrument-df/
type CandleStickData struct {
	Open  string `json:"o"`
	High  string `json:"h"`
	Low   string `json:"l"`
	Close string `json:"c"`
}
