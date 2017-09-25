package main

import (
	"fmt"
	"os"

	"github.com/kkastan/goanda/candles"
	"github.com/kkastan/goanda/ticker"
)

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
	}

	fmt.Printf("Unrecognized option: %s\n", args[0])
	os.Exit(-1)
}

func candleClient() {
	data := candles.Get()

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
