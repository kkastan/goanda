package main

import (
	"fmt"

	"github.com/kkastan/goanda/ticker"
)

func main() {
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
