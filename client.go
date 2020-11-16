package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/kkastan/goanda/candles"
	"github.com/kkastan/goanda/common"
	"github.com/kkastan/goanda/order"
	"github.com/kkastan/goanda/ticker"
)

type logger struct {
}

func (l *logger) Infof(format string, args ...interface{}) {
	fmt.Printf(format+"\n", args...)
}

func main() {

	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Printf("Usage: client [ticker|candles]")
		os.Exit(-1)
	}

	switch {
	case args[0] == "candles":
		candleClient()
		return
	case args[0] == "ticker":
		tickerClient()
		return
	case args[0] == "market":
		marketOrderClient()
		return
	}

	fmt.Printf("Unrecognized option: %s\n", args[0])
	os.Exit(-1)
}

func marketOrderClient() {
	l := &logger{}

	marketCmd := flag.NewFlagSet("marketOrder", flag.ExitOnError)
	stopLoss := marketCmd.Float64("s", 0.0, "Stop Loss")
	takeProfit := marketCmd.Float64("p", 0.0, "Take Profit")
	quantity := marketCmd.Int64("q", 0, "Quantity - negative for short")

	marketCmd.Parse(os.Args[2:])

	l.Infof("Stop Loss set at %f", *stopLoss)

	ch := make(chan *ticker.Tick)
	t := ticker.New("EUR_USD", ch)
	t.Run()

	select {
	case tick := <-ch:
		l.Infof("%s\t%s\tBID/ASK\t%f/%f %f", tick.Time, tick.Instrument, tick.GetBid(), tick.GetAsk(), tick.GetAsk()-tick.GetBid())

		sl := tick.GetAsk() + *stopLoss
		tp := tick.GetAsk() + *takeProfit

		l.Infof("Executing order: %d units at %f with stop loss=%f and take profit=%f", *quantity, tick.GetAsk(), sl, tp)

		o := &order.Orderer{
			Log: l,
		}

		orderRequest := &order.Order{
			Type:        order.MARKET,
			Units:       *quantity,
			Instrument:  "EUR_USD",
			TimeInForce: order.FOK,
			StopLoss:    sl,
			TakeProfit:  tp,
		}

		o.SubmitOrder(orderRequest)
	}

}

func candleClient() {
	fmt.Printf("Recent Candles...\n")
	count := int32(10)

	cr := &candles.CandleRequest{
		Price:        "M",
		Granularity:  "M1",
		Count:        &count,
		IncludeFirst: true,
	}

	data := candles.Get("EUR_USD", cr)

	fmt.Printf("%s\t%s\n", data.Instrument, data.Granularity)

	for _, c := range data.Candles {
		fmt.Printf("%s\t%t\t%d\t%s\t%s\t%s\t%s\n", c.Time, c.Complete, c.Volume,
			c.Mid.Open, c.Mid.High, c.Mid.Low, c.Mid.Close)
	}

	from, err := time.Parse(common.OandaRFC3339Format, "2017-09-20T00:00:00.000000000Z+00:00")
	if err != nil {
		panic(err.Error())
	}

	to, err := time.Parse(common.OandaRFC3339Format, "2017-09-21T00:00:00.000000000Z+00:00")
	if err != nil {
		panic(err.Error())
	}

	cr = &candles.CandleRequest{
		Price:        "M",
		Granularity:  "M15",
		From:         from,
		To:           to,
		IncludeFirst: true,
	}

	fmt.Printf("\nTime Range Candles...\n")
	data = candles.Get("EUR_USD", cr)
	fmt.Printf("%s\t%s\n", data.Instrument, data.Granularity)

	for _, c := range data.Candles {
		fmt.Printf("%s\t%t\t%d\t%s\t%s\t%s\t%s\n", c.Time, c.Complete, c.Volume,
			c.Mid.Open, c.Mid.High, c.Mid.Low, c.Mid.Close)
	}
}

func tickerClient() {
	ch := make(chan *ticker.Tick)
	t := ticker.New("EUR_USD,GBP_USD", ch)
	t.Run()

	for {
		select {
		case tick := <-ch:
			fmt.Printf("%s\t%s\tBID/ASK\t%f/%f\n", tick.Time, tick.Instrument, tick.GetBid(), tick.GetAsk())
		}
	}
}
