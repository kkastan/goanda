package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/kkastan/goanda/candles"
	"github.com/kkastan/goanda/common"
	"github.com/kkastan/goanda/order"
	"github.com/kkastan/goanda/positions"
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
	case args[0] == "latest":
		latestCandleClient()
		return
	case args[0] == "positions":
		openPositionsClient()
		return
	case args[0] == "order":
		price, err := strconv.ParseFloat(args[1], 64)
		if err != nil {
			panic(err)
		}
		orderClient(price)
		return
	}

	fmt.Printf("Unrecognized option: %s\n", args[0])
	os.Exit(-1)
}

func orderClient(price float64) {

	log := &logger{}
	o := &order.Orderer{
		Log: log,
	}

	log.Infof("orderClient called with price: %f", price)

	var quantity int64 = 100
	stop := price * 0.9995
	profit := price * 1.001

	log.Infof("stop loss:\t%f", stop)
	log.Infof("take profit:\t%f", profit)

	orderRequest := &order.Order{
		Type:        order.MARKET,
		Units:       quantity,
		Instrument:  "GBP_USD",
		TimeInForce: order.FOK,
		StopLoss:    stop,
		TakeProfit:  profit,
		ClientID:    "client.id.1001",
		ClientTag:   "clinet.tag.1002",
	}

	log.Infof("%v", orderRequest)

	resp, err := o.SubmitOrder(orderRequest)
	if err != nil {
		log.Infof("error: %v\n", err)
	}

	log.Infof("%v", resp)
}

func openPositionsClient() {
	l := &logger{}
	client := &positions.Client{
		Log: l,
	}

	results, err := client.GetOpenPositions()
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(results)
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

		resp, err := o.SubmitOrder(orderRequest)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(resp)
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
	t := ticker.New("EUR_USD,GBP_USD,EUR_GBP", ch)
	t.Run()

	for {
		select {
		case tick := <-ch:
			fmt.Printf("%s\t%s\tBID/ASK\t%f/%f\n", tick.Time, tick.Instrument, tick.GetBid(), tick.GetAsk())
		}
	}
}

func latestCandleClient() {
	data := candles.GetLatest("EUR_USD:M5:BMA,EUR_GBP:M5:BMA,GBP_USD:M5:BMA")
	str, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", str)
}
