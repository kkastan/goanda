package candles

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/kkastan/goanda/common"
)

// CandleRequest encapsulates the various options in
// requesting a set of candles. See the candles endpoint
// documentation at http://developer.oanda.com/rest-live-v20/instrument-ep/
type CandleRequest struct {
	Price             string
	Granularity       string
	Count             *int32
	From              time.Time
	To                time.Time
	Smooth            bool // Oanda defaults to False
	IncludeFirst      bool // Oanda defaults to True
	DailyAlignment    *int32
	AlignmentTimezone string
	WeeklyAlignment   string
}

func parseTimeAsOandaRFC3339String(t time.Time) (str string) {
	t = t.In(time.UTC)
	str = t.Format(common.OandaRFC3339Format)
	// remove the trailing +00:00
	str = str[:len(str)-6]
	return
}

func constructCandleParams(cr *CandleRequest) (params string) {
	params = fmt.Sprintf("?price=%s&granularity=%s", cr.Price, cr.Granularity)

	if cr.Count != nil {
		params = fmt.Sprintf("%s&count=%d", params, *cr.Count)
	} else {
		// convert times
		from := parseTimeAsOandaRFC3339String(cr.From)
		to := parseTimeAsOandaRFC3339String(cr.To)
		params = fmt.Sprintf("%s&from=%s&to=%s", params, from, to)

		// IncludeFirst only applies if the from parameter is used
		if cr.IncludeFirst {
			params = fmt.Sprintf("%s&includeFirst=True", params)
		} else {
			params = fmt.Sprintf("%s&includeFirst=False", params)
		}
	}

	if cr.Smooth {
		params = fmt.Sprintf("%s&smooth=True", params)
	} else {
		params = fmt.Sprintf("%s&smooth=False", params)
	}

	if cr.DailyAlignment != nil {
		params = fmt.Sprintf("%s&dailyAlignment=%d", params, cr.DailyAlignment)
	}

	if cr.AlignmentTimezone != "" {
		params = fmt.Sprintf("%s&alignmentTimezone=%s", params, cr.AlignmentTimezone)
	}

	if cr.WeeklyAlignment != "" {
		params = fmt.Sprintf("%s&weeklyAlignment=%s", params, cr.WeeklyAlignment)
	}

	return
}

// Get is a simple wrapper of the Oanda v20 candles endpoint -
// see http://developer.oanda.com/rest-live-v20/instrument-ep/
func Get(instrument string, cr *CandleRequest) *CandleResponse {
	baseURL := os.Getenv(common.FxAPIBaseURLEnvVarName)
	url := fmt.Sprintf("%s/instruments/%s/candles%s",
		baseURL, instrument, constructCandleParams(cr))

	log.Printf("url: %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(fmt.Errorf("error creating the request %s", err.Error()))
	}

	bearerToken := os.Getenv("OANDA_API_KEY")
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
