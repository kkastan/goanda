package order

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_constructPayloadFromRequest(t *testing.T) {
	oo := &Order{
		Type:        MARKET,
		Units:       10000,
		Instrument:  "EUR_USD",
		TimeInForce: FOK,
		StopLoss:    1.17432,
		TakeProfit:  1.17452,
	}

	p, err := constructPayloadFromRequest(oo)

	if assert.Nil(t, err) {
		if assert.NotNil(t, p) {
			if assert.NotNil(t, p.Order) {
				assert.Equal(t, Type(MARKET), p.Order.Type)
				assert.Equal(t, "10000", p.Order.Units)
				assert.Equal(t, "EUR_USD", p.Order.Instrument)
				assert.Equal(t, TimeInForce(FOK), p.Order.TimeInForce)
				if assert.NotNil(t, p.Order.StopLossOnFill) {
					assert.Equal(t, "1.17432", p.Order.StopLossOnFill.Price)
					assert.Equal(t, TimeInForce(GTC), p.Order.StopLossOnFill.TimeInForce)
				}
				if assert.NotNil(t, p.Order.TakeProfitOnFill) {
					assert.Equal(t, "1.17452", p.Order.TakeProfitOnFill.Price)
					assert.Equal(t, TimeInForce(GTC), p.Order.TakeProfitOnFill.TimeInForce)
				}
			}
		}
	}
}

func Test_constructPayloadFromRequest_ValidStopLossTakeProfitLong(t *testing.T) {
	oo := &Order{
		Type:        MARKET,
		Units:       10000,
		Instrument:  "EUR_USD",
		TimeInForce: FOK,
		StopLoss:    1.17452,
		TakeProfit:  1.17432,
	}

	p, err := constructPayloadFromRequest(oo)

	if assert.NotNil(t, err) {
		assert.Equal(t, "stop loss must be less than the take profit on a long position", err.Error())
		assert.Nil(t, p)
	}

	oo.Units = -10000
	oo.TakeProfit = 1.17472
	oo.StopLoss = 1.17462

	p, err = constructPayloadFromRequest(oo)

	if assert.NotNil(t, err) {
		assert.Equal(t, "stop loss must be greater than the take profit on a short position", err.Error())
		assert.Nil(t, p)
	}
}

func Test_constructPayloadFromRequest_MarketTimeInForce(t *testing.T) {

	TIFs := []TimeInForce{GTC, GTD, GFD}

	for _, tif := range TIFs {

		oo := &Order{
			Type:        MARKET,
			Units:       10000,
			Instrument:  "EUR_USD",
			TimeInForce: tif,
			StopLoss:    1.17432,
			TakeProfit:  1.17452,
		}

		p, err := constructPayloadFromRequest(oo)

		if assert.NotNil(t, err) {
			assert.Equal(t, "market orders can only be FOK or IOC", err.Error())
			assert.Nil(t, p)
		}
	}
}
