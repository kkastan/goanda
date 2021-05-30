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

type logger interface {
	Infof(format string, args ...interface{})
}

// Orderer ...
type Orderer struct {
	Log logger
}

// OandaOrderResponse ...
type OandaOrderResponse map[string]interface{}

// SubmitOrder ...
func (o *Orderer) SubmitOrder(order *Order) (oresp *OandaOrderResponse, err error) {
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

	o.Log.Infof("posting order. url: %s\nbody: %s", url, string(raw))

	// fmt.Printf("url: %s\n", url)
	// fmt.Printf("output:\n%s\n", string(raw))

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

	o.Log.Infof("response Status: %s", resp.Status)
	o.Log.Infof("response Headers: %s", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	o.Log.Infof("response Body: %s", string(body))

	oresp = &OandaOrderResponse{}
	err = json.Unmarshal(body, oresp)

	if resp.StatusCode != 201 {
		if errorMessage, ok := (*oresp)["errorMessage"]; ok {
			err = fmt.Errorf(errorMessage.(string))
		} else {
			err = fmt.Errorf("Oanda orders endpoint returned status %d expected 201", resp.StatusCode)
		}
	}

	return
}

func constructPayloadFromRequest(order *Order) (oop *oandaOrderPayload, err error) {

	if order.Type == MARKET && !(order.TimeInForce == FOK || order.TimeInForce == IOC) {
		err = fmt.Errorf("market orders can only be FOK or IOC")
		return
	}

	if order.Units > 0 &&
		order.StopLoss > 0 &&
		order.TakeProfit > 0 &&
		order.StopLoss >= order.TakeProfit {
		err = fmt.Errorf("stop loss must be less than the take profit on a long position")
		return
	}

	if order.Units < 0 &&
		order.StopLoss > 0 &&
		order.TakeProfit > 0 &&
		order.StopLoss <= order.TakeProfit {
		err = fmt.Errorf("stop loss must be greater than the take profit on a short position")
		return
	}

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
