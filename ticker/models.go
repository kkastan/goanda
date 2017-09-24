package ticker

import "strconv"

// Tick data struction
type Tick struct {
	Time string `json:"time,omitempty"`
	Type string `json:"type"`
	// The v20 API exposes Oanda's depth of market in its pricing
	// data. See the answer to "What is Depth of Market (DOM)?"
	// at https://oanda.secure.force.com/AnswersSupport?language=en_US&urlName=V20-Trading-Platform-FAQ
	// for more details. For our purposes we'll be assuming that
	// we're trading at the lower end of the liquidity specturm
	// and will take the highest bid price and lowest ask price
	// as the "effective" bid and ask. See GetBid() and GetAsk().
	// If you're trading at the level where you need to account
	// for the multiple bid and ask prices as your volume per
	// trade spans liquidity depths then more power to you. That's
	// awesome and you should be able afford the time to extend
	// this code to account for that or hire some engineers
	// to do so for you.
	Bids       []Quote `json:"bids"`
	Asks       []Quote `json:"asks"`
	Instrument string  `json:"instrument"`
	Tradeable  bool    `json:"tradeable"`
}

// Quote data structure
type Quote struct {
	Price     string `json:"price"`
	Liquidity int32  `json:"liquidity"`
}

// GetBid gets the highest bid. See the Bids/Asks comment
// above in the Tick struct for more information.
func (t *Tick) GetBid() (bid float64) {
	if 0 == len(t.Bids) {
		return
	}

	bid = t.Bids[0].PriceAsFloat()

	// best bid is the highest?
	for _, quote := range t.Bids {
		if quote.PriceAsFloat() > bid {
			bid = quote.PriceAsFloat()
		}
	}

	return
}

// GetAsk gets the lowest ask. See the Bids/Asks comment
// above in the Tick struct for more information.
func (t *Tick) GetAsk() (ask float64) {
	if 0 == len(t.Asks) {
		return
	}

	ask = t.Asks[0].PriceAsFloat()

	for _, quote := range t.Asks {
		if quote.PriceAsFloat() < ask {
			ask = quote.PriceAsFloat()
		}
	}

	return
}

// PriceAsFloat converts the price to a float64
func (q *Quote) PriceAsFloat() (val float64) {
	val, err := strconv.ParseFloat(q.Price, 64)
	if err != nil {
		panic(err.Error())
	}

	return
}
