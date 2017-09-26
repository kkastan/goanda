package candles

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/kkastan/goanda/common"
)

// GetRecent returns the last <i>count</i> candles.
// Simple wrapper of the Oanda v20 candles endpoint -
// see http://developer.oanda.com/rest-live-v20/instrument-ep/
func GetRecent(instrument, price, granularity string, count int32) *CandleResponse {
	bearerToken := os.Getenv("OANDA_API_KEY")
	baseURL := os.Getenv(common.FxAPIBaseURLEnvVarName)

	url := fmt.Sprintf("%s/instruments/%s/candles?count=%d&price=%s&granularity=%s",
		baseURL, instrument, count, price, granularity)

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
// Simple time range wrapper of the Oanda v20 candles endpoint -
// see http://developer.oanda.com/rest-live-v20/instrument-ep/
func GetTimeRange(instrument, price, granularity, from, to string) *CandleResponse {
	bearerToken := os.Getenv("OANDA_API_KEY")
	baseURL := os.Getenv(common.FxAPIBaseURLEnvVarName)

	url := fmt.Sprintf("%s/instruments/%s/candles?price=%s&granularity=%s&from=%s&to=%s&includeFirst=true",
		baseURL, instrument, price, granularity, from, to)

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
