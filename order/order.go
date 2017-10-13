package order

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/kkastan/goanda/common"
)

// SubmitOrder ...
func SubmitOrder(order *Order) (err error) {
	p, err := constructPayloadFromRequest(order)
	if err != nil {
		return
	}

	raw, err := json.Marshal(p)
	if err != nil {
		return
	}

	baseURL := os.Getenv(common.FxAPIBaseURLEnvVarName)
	accountID := os.Getenv(common.AccountIDEnvVarName)
	bearerToken := os.Getenv(common.APIKeyEnvVarName)

	url := fmt.Sprintf("%s/accounts/%s/orders", baseURL, accountID)

	fmt.Printf("url: %s\n", url)
	fmt.Printf("output:\n%s\n", string(raw))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(raw))
	if err != nil {
		panic(fmt.Errorf("error creating the request %s", err.Error()))
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Accept-Datetime-Format", "RFC3339")
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Errorf("GET error: %s", err.Error()))
	}

	defer resp.Body.Close()

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))

	return
}

// TODO unit tests - lost of them
// test that the stop order and take profit make sense for a long/short order
// review the Order API and validate that the requested set of options are allowed
func constructPayloadFromRequest(order *Order) (oop *oandaOrderPayload, err error) {
	oo := &oandaOrder{}

	oo.Type = order.Type
	oo.Units = strconv.FormatInt(order.Units, 10)
	oo.Instrument = order.Instrument
	oo.PositionFill = DEFAULT

	if order.TimeInForce == "" {
		oo.TimeInForce = FOK
	} else {
		oo.TimeInForce = order.TimeInForce
	}

	if order.StopLoss != 0 {
		oo.StopLossOnFill = &oandaStopLossDetails{
			Price:       strconv.FormatFloat(order.StopLoss, 'f', 5, 64),
			TimeInForce: GTC,
		}
	}

	if order.TakeProfit != 0 {
		oo.TakeProfitOnFill = &oandaTakeProfitDetails{
			Price:       strconv.FormatFloat(order.TakeProfit, 'f', 5, 64),
			TimeInForce: GTC,
		}
	}

	oop = &oandaOrderPayload{
		Order: oo,
	}

	return
}
