package ticker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

const (
	apiKeyEnvVarName           = "OANDA_API_KEY"
	accountIDEnvVarName        = "OANDA_ACCOUNT_ID"
	streamAPIBaseURLEnvVarName = "OANDA_STREAM_API"
)

// Ticker data structure
type Ticker struct {
	Instruments string
	Ticks       chan *Tick
}

// New constructs a new Ticker
func New(instruments string, c chan *Tick) *Ticker {
	return &Ticker{
		Instruments: instruments,
		Ticks:       c,
	}
}

// Run runs the ticker in a goroutine
func (t *Ticker) Run() {
	go t.runInternal()
}

func (t *Ticker) runInternal() {

	baseURL := os.Getenv(streamAPIBaseURLEnvVarName)
	accountID := os.Getenv(accountIDEnvVarName)
	bearerToken := os.Getenv(apiKeyEnvVarName)

	if baseURL == "" || accountID == "" || bearerToken == "" {
		panic(fmt.Sprintf("One or more of the following environment variables "+
			"is empty or not set:\n%s\n%s\n%s",
			streamAPIBaseURLEnvVarName, accountIDEnvVarName, apiKeyEnvVarName))
	}

	url := t.getPriceStreamURL(baseURL, accountID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(fmt.Errorf("error creating the request %s", err.Error()))
	}

	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(fmt.Errorf("GET error: %s", err.Error()))
	}

	reader := bufio.NewReader(resp.Body)

	tick := &Tick{}

	for {
		line, err := reader.ReadBytes('\n')
		if err != nil {
			panic(fmt.Errorf("reader.ReadBytes: %s", err.Error()))
		}

		if err := json.Unmarshal(line, tick); err != nil {
			panic(fmt.Errorf("json.Unmarshal: %s", err.Error()))
		}

		if tick.Type == "PRICE" && tick.Tradeable {
			t.Ticks <- tick
		}

	}
}

func (t *Ticker) getPriceStreamURL(baseURL, accountID string) (url string) {
	url = fmt.Sprintf("%s/accounts/%s/pricing/stream?instruments=%s", baseURL, accountID, t.Instruments)
	return
}
