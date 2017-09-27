package main

import (
	"fmt"
	"os"
	"time"

	"github.com/kkastan/goanda/candles"
	"github.com/kkastan/goanda/common"
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

	from, err := time.Parse(common.OandaRFC3339Format, "2017-09-20T00:00:00.00000000Z+00:00")
	if err != nil {
		panic(err.Error())
	}

	to, err := time.Parse(common.OandaRFC3339Format, "2017-09-21T00:00:00.00000000Z+00:00")
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
