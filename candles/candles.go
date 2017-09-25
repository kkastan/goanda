package candles

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func Get() *CandleResponse {

	bearerToken := os.Getenv("OANDA_API_KEY")
	url := "https://api-fxpractice.oanda.com/v3/instruments/EUR_USD/candles?count=6&price=M&granularity=M1"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(fmt.Errorf("error creating the request %s", err.Error()))
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Errorf("GET error: %s", err.Error()))
	}

	candleResponse := &CandleResponse{}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(candleResponse)

	return candleResponse
}
