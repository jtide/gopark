package api_test

import (
	"github.com/jtide/gopark/api"
	"github.com/stretchr/testify/assert"
	"testing"
)

var jsonStandardConfig = []byte(
	`{
    "rates": [
        {
            "days": "mon,tues,thurs",
            "times": "0900-2100",
            "price": 1500
        },
        {
            "days": "fri,sat,sun",
            "times": "0900-2100",
            "price": 2000
        },
        {
            "days": "wed",
            "times": "0600-1800",
            "price": 1750
        },
        {
            "days": "mon,wed,sat",
            "times": "0100-0500",
            "price": 1000
        },
        {
            "days": "sun,tues",
            "times": "0100-0700",
            "price": 925
        }
    ]
}`)

func TestUpdateRatesWithStandardConfig(t *testing.T) {
	rates := api.NewWeeklyRates()
	err := rates.Update(jsonStandardConfig)
	assert.NoError(t, err)
}

func TestUpdateRatesWithValidConfig(t *testing.T) {
	json := []byte(
		`{
  "rates": [
	{
	  "days": "mon,tues,wed,thurs,fri",
	  "times": "0600-1800",
	  "price": 1500
	},
	{
	  "days": "sat,sun",
	  "times": "0600-2000",
	  "price": 2000
	}
  ]
}`)
	rates := api.NewWeeklyRates()
	err := rates.Update(json)
	assert.NoError(t, err)
}

func TestUpdateRatesWithInvalidDay(t *testing.T) {
	json := []byte(
		`{
  "rates": [
	{
	  "days": "mon,tues,wed,thurs,InvalidWeekday",
	  "times": "0600-1800",
	  "price": 1500
	},
	{
	  "days": "sat,sun",
	  "times": "0600-2000",
	  "price": 2000
	}
  ]
}`)
	rates := api.NewWeeklyRates()
	err := rates.Update(json)
	assert.Error(t, err, "could not update rate: 'InvalidWeekday' is not a recognized weekday")
}

func TestUpdateRatesWithInvalidMinutes(t *testing.T) {
	json := []byte(
		`{
  "rates": [
	{
	  "days": "mon,tues,wed,thurs,fri",
	  "times": "0600-1860",
	  "price": 1500
	},
	{
	  "days": "sat,sun",
	  "times": "0600-2000",
	  "price": 2000
	}
  ]
}`)
	rates := api.NewWeeklyRates()
	err := rates.Update(json)
	assert.Error(t, err, "could not update rate: invalid end time: invalid minutes in time: 1860")
}

func TestUpdateRatesWithOverlappingTimes1(t *testing.T) {
	rates := api.NewWeeklyRates()
	json := []byte(
		`{
  "rates": [
	{
	  "days": "mon,tues,wed,thurs,fri",
	  "times": "0600-1800",
	  "price": 1500
	},
	{
	  "days": "mon",
	  "times": "0601-0700",
	  "price": 1500
	}
  ]
}`)
	err := rates.Update(json)
	assert.Errorf(t, err, "could not update rate: new rate presents a conflict: {Monday 361 420 1500}")
}

func TestUpdateRatesWithOverlappingTimes2(t *testing.T) {
	json := []byte(
		`{
  "rates": [
	{
	  "days": "mon,tues,wed,thurs,fri",
	  "times": "0600-1800",
	  "price": 1500
	},
	{
	  "days": "mon",
	  "times": "0501-0601",
	  "price": 1500
	}
  ]
}`)
	rates := api.NewWeeklyRates()
	err := rates.Update(json)
	assert.Errorf(t, err, "could not update rate: new rate presents a conflict {Monday 301 361 1500}: new rate conficts with existing rate {Monday 360 1080 1500}")
}

func TestLookupPriceFromDifferentDays(t *testing.T) {
	rates := api.NewWeeklyRates()
	err := rates.Update(jsonStandardConfig)
	assert.NoError(t, err)

	price, err := rates.Lookup("2018-05-02T06:00:00Z", "2018-05-09T16:30:00Z")
	assert.Errorf(t, err, "start and end times must be on same day: start=122, end=129")
	assert.Equal(t, uint(0), price)
}

func TestLookupPriceAvailable1(t *testing.T) {
	rates := api.NewWeeklyRates()
	err := rates.Update(jsonStandardConfig)
	assert.NoError(t, err)

	price, err := rates.Lookup("2015-07-01T07:00:00Z", "2015-07-01T16:00:00Z")
	assert.NoError(t, err)
	assert.Equal(t, uint(1750), price)
}

func TestLookupPriceAvailable2(t *testing.T) {
	rates := api.NewWeeklyRates()
	err := rates.Update(jsonStandardConfig)
	assert.NoError(t, err)

	price, err := rates.Lookup("2018-05-02T06:00:00Z", "2018-05-02T16:30:00Z")
	assert.NoError(t, err)
	assert.Equal(t, uint(1750), price)
}

func TestLookupPriceUnavailable1(t *testing.T) {
	rates := api.NewWeeklyRates()
	err := rates.Update(jsonStandardConfig)
	assert.NoError(t, err)

	price, err := rates.Lookup("2018-05-01T06:00:00Z", "2018-05-01T16:30:00Z")
	assert.Error(t, err)
	assert.Equal(t, uint(0), price)
}

func TestLookupPriceUnavailable2(t *testing.T) {
	rates := api.NewWeeklyRates()
	err := rates.Update(jsonStandardConfig)
	assert.NoError(t, err)

	price, err := rates.Lookup("2018-05-02T05:59:59Z", "2018-05-02T16:30:00Z")
	assert.Error(t, err)
	assert.Equal(t, uint(0), price)
}

func TestWeeklyRates_DeepCopy(t *testing.T) {
	rates := api.NewWeeklyRates()
	err := rates.Update(jsonStandardConfig)
	assert.NoError(t, err)

	ratesCopy := rates.DeepCopy()

	var jsonUpdate = []byte(
		`{
    "rates": [
        {
            "days": "mon",
            "times": "0600-0700",
            "price": 1234
        }
    ]
}`)

	// Only update the copy
	err = ratesCopy.Update(jsonUpdate)
	assert.NoError(t, err)

	// HourlyRate from original should be found in copy
	price, err := ratesCopy.Lookup("2018-04-30T04:30:00Z", "2018-04-30T04:45:00Z")
	assert.NoError(t, err)
	assert.Equal(t, uint(1000), price)

	// HourlyRate from update should be found in copy
	price, err = ratesCopy.Lookup("2018-04-30T06:00:00Z", "2018-04-30T06:30:00Z")
	assert.NoError(t, err)
	assert.Equal(t, uint(1234), price)

	// HourlyRate from update should not be found in original
	price, err = rates.Lookup("2018-04-30T06:00:00Z", "2018-04-30T06:30:00Z")
	assert.Error(t, err)
	assert.Equal(t, uint(0), price)
}
