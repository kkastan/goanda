package ticker

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kkastan/goanda/common"
)

const (
	MAX_ERROR_COUNT = 5
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

	baseURL := os.Getenv(common.StreamAPIBaseURLEnvVarName)
	accountID := os.Getenv(common.AccountIDEnvVarName)
	bearerToken := os.Getenv(common.APIKeyEnvVarName)

	if baseURL == "" || accountID == "" || bearerToken == "" {
		panic(fmt.Sprintf("One or more of the following environment variables "+
			"is empty or not set:\n%s\n%s\n%s",
			common.StreamAPIBaseURLEnvVarName, common.AccountIDEnvVarName,
			common.APIKeyEnvVarName))
	}

	url := t.getPriceStreamURL(baseURL, accountID)

	// Outer infinte loop allows us to re-establish the connection to Oanda
	// if a threshold of errors are exceeded when reading from the server.
	for {
		log.Println("Establishing connection to Oanda...")

		errorCount := 0
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

		// Infinte loop to read each new line as sent from the server. If too many
		// errors occur then break out of this loop to have the connection
		// re-established.
		for {
			if errorCount > MAX_ERROR_COUNT {
				log.Printf("Exceeded %d errors in reading from Oanda. Re-establishing connection to Oanda.", errorCount)
				break
			}

			line, err := reader.ReadBytes('\n')
			if err != nil {
				log.Printf("reader.ReadBytes [%v]: %s", line, err.Error())
				errorCount++
				continue
			}

			// fmt.Printf("%s\n", line)

			if err := json.Unmarshal(line, tick); err != nil {
				log.Printf("json.Unmarshal: %s", err.Error())
				errorCount++
				continue
			}

			if tick.Type == "PRICE" && tick.Tradeable {
				t.Ticks <- tick
			}

		} // end inner infinte loop
	} // end outer infinte loop
}

func (t *Ticker) getPriceStreamURL(baseURL, accountID string) (url string) {
	url = fmt.Sprintf("%s/accounts/%s/pricing/stream?instruments=%s", baseURL, accountID, t.Instruments)
	return
}
