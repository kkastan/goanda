package candles

import (
	"regexp"
	"testing"
	"time"

	"github.com/kkastan/goanda/common"
	"github.com/stretchr/testify/assert"
)

func Test_parseTimeAsOandaRFC3339String(t *testing.T) {

	testTime, err := time.Parse("Jan 2 15:04:05.99999999 2006-07:00", "Sep 20 13:00:00.12345 2017-07:00")
	if assert.Nil(t, err) {
		str := parseTimeAsOandaRFC3339String(testTime)
		assert.Equal(t, "2017-09-20T20:00:00.123450000Z", str)
	}

	testTime, err = time.Parse(common.OandaRFC3339Format, "2017-10-01T04:15:03.123400000Z+00:00")
	if assert.Nil(t, err) {
		str := parseTimeAsOandaRFC3339String(testTime)
		assert.Equal(t, "2017-10-01T04:15:03.123400000Z", str)
	}
}

func Test_constructCandleParamsTimeRange(t *testing.T) {

	from, err := time.Parse(common.OandaRFC3339Format, "2017-09-20T00:00:00.000000000Z+00:00")
	if !assert.Nil(t, err, "Error parsing from time") {
		return
	}

	to, err := time.Parse(common.OandaRFC3339Format, "2017-09-21T00:00:00.000000000Z+00:00")
	if !assert.Nil(t, err, "Error parsing to time") {
		return
	}

	cr := &CandleRequest{
		Price:       "M",
		Granularity: "M15",
		From:        from,
		To:          to,
	}

	params := constructCandleParams(cr)

	m, err := regexp.MatchString("price=M", params)
	if assert.Nil(t, err) {
		assert.True(t, m, "Unexpected price")
	}

	m, err = regexp.MatchString("granularity=M15", params)
	if assert.Nil(t, err) {
		assert.True(t, m, "Unexpected granularity")
	}

	m, err = regexp.MatchString("from=2017-09-20T00:00:00.000000000Z", params)
	if assert.Nil(t, err) {
		assert.True(t, m, "Unexpected from time")
	}

	m, err = regexp.MatchString("to=2017-09-21T00:00:00.000000000Z", params)
	if assert.Nil(t, err) {
		assert.True(t, m, "Unexpected to time")
	}
}

func Test_constructCandleParamsSimpleCount(t *testing.T) {
	count := int32(10)
	cr := &CandleRequest{
		Price:       "M",
		Granularity: "M5",
		Count:       &count,
	}
	params := constructCandleParams(cr)

	m, err := regexp.MatchString("price=M", params)
	if assert.Nil(t, err) {
		assert.True(t, m, "Unexpected price")
	}

	m, err = regexp.MatchString("granularity=M5", params)
	if assert.Nil(t, err) {
		assert.True(t, m, "Unexpected granularity")
	}

	m, err = regexp.MatchString("count=10", params)
	if assert.Nil(t, err) {
		assert.True(t, m, "Unexpected count")
	}

	m, err = regexp.MatchString("to=", params)
	if assert.Nil(t, err) {
		assert.False(t, m, "Unexpected to param")
	}

	m, err = regexp.MatchString("from=", params)
	if assert.Nil(t, err) {
		assert.False(t, m, "Unexpected from param")
	}

	m, err = regexp.MatchString("smooth=False", params)
	if assert.Nil(t, err) {
		assert.True(t, m, "Unexpected smooth")
	}

	m, err = regexp.MatchString("includeFirst=", params)
	if assert.Nil(t, err) {
		assert.False(t, m, "Unexpected includeFirst")
	}

	m, err = regexp.MatchString("dailyAlignment=", params)
	if assert.Nil(t, err) {
		assert.False(t, m, "Unexpected dailyAlignment param")
	}

	m, err = regexp.MatchString("alignmentTimezone=", params)
	if assert.Nil(t, err) {
		assert.False(t, m, "Unexpected alignmentTimezone param")
	}

	m, err = regexp.MatchString("weeklyAlignment=", params)
	if assert.Nil(t, err) {
		assert.False(t, m, "Unexpected weeklyAlignment param")
	}

}
