package main

//
// import (
// 	"fmt"
//
// 	"github.com/kkastan/goanda/order"
// )
//
// type logger struct {
// }
//
// func (l *logger) Infof(format string, args ...interface{}) {
// 	fmt.Printf(format+"\n", args...)
// }
//
// func main() {
//
// 	var quantity int64 = 100
// 	price := 1.21921
// 	stop := price * 0.9995
// 	profit := price * 1.001
//
// 	fmt.Println(stop)
// 	fmt.Println(profit)
//
// 	orderRequest := &order.Order{
// 		Type:        order.MARKET,
// 		Units:       quantity,
// 		Instrument:  "EUR_USD",
// 		TimeInForce: order.FOK,
// 		StopLoss:    stop,
// 		TakeProfit:  profit,
// 	}
//
// 	fmt.Println(orderRequest)
//
// 	l := &logger{}
// 	o := &order.Orderer{
// 		Log: l,
// 	}
//
// 	resp, err := o.SubmitOrder(orderRequest)
// 	if err != nil {
// 		fmt.Printf("error: %v\n", err)
// 	} else {
// 		fmt.Println("no error!")
// 	}
// 	fmt.Println(resp)
// }
