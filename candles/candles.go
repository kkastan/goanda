package candles

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// GetRecent returns the last N candles
func GetRecent() *CandleResponse {
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

// GetTimeRange returns the candles for the specified time range
func GetTimeRange() *CandleResponse {
	bearerToken := os.Getenv("OANDA_API_KEY")
	//	url := "https://api-fxpractice.oanda.com/v3/instruments/EUR_USD/candles?price=M&granularity=M5"

	// YYYY-MM-DDTHH:MM:SS.nnnnnnnnnZ
	from := "2017-09-20T00:00:00.0Z"
	to := "2017-09-21T00:00:00.0Z"
	url := fmt.Sprintf("https://api-fxpractice.oanda.com/v3/instruments/%s/candles?price=%s&granularity=%s&from=%s&to=%s&includeFirst=true",
		"EUR_USD", "M", "M5", from, to)

	fmt.Printf("url: %s\n", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(fmt.Errorf("error creating the request %s", err.Error()))
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)
	req.Header.Set("Accept-Datetime-Format", "RFC3339")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Errorf("GET error: %s", err.Error()))
	}

	candleResponse := &CandleResponse{}
	defer resp.Body.Close()
	json.NewDecoder(resp.Body).Decode(candleResponse)

	return candleResponse
}
