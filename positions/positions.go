package positions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/kkastan/goanda/common"
)

type logger interface {
	Infof(format string, args ...interface{})
}

// Client ...
type Client struct {
	Log logger
}

// GetOpenPositions see https://developer.oanda.com/rest-live-v20/position-ep/
func (c *Client) GetOpenPositions() (oopr *OandaOpenPositionsResults, err error) {
	baseURL := os.Getenv(common.FxAPIBaseURLEnvVarName)
	accountID := os.Getenv(common.AccountIDEnvVarName)

	url := fmt.Sprintf("%s/accounts/%s/openPositions",
		baseURL, accountID)

	c.Log.Infof("url: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	bearerToken := os.Getenv("OANDA_API_KEY")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	c.Log.Infof("response Body: %s", string(body))

	oopr = &OandaOpenPositionsResults{}
	err = json.Unmarshal(body, oopr)

	return
}
