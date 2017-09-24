package ticker

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

const tick1 = `{"type":"PRICE","time":"2017-09-24T22:20:39.888097022Z","bids":[{"price":"1.19195","liquidity":10000000},{"price":"1.19197","liquidity":500000}],"asks":[{"price":"1.19218","liquidity":10000000},{"price":"1.19216","liquidity":500000}],"closeoutBid":"1.19180","closeoutAsk":"1.19233","status":"tradeable","tradeable":true,"instrument":"EUR_USD"}`

func Test_GetBid_GetsHighestBid(t *testing.T) {
	tick := &Tick{}
	err := json.Unmarshal([]byte(tick1), tick)
	if assert.Nil(t, err, "error unmarshalling test tick data") {
		assert.Equal(t, 1.19197, tick.GetBid())
	}
}

func Test_GetAsk_GetsLowestAsk(t *testing.T) {
	tick := &Tick{}
	err := json.Unmarshal([]byte(tick1), tick)
	if assert.Nil(t, err, "error unmarshalling test tick data") {
		assert.Equal(t, 1.19216, tick.GetAsk())
	}

}
