package api

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"
)

type ReturnRate struct {
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Price uint      `json:"price"`
}

// JSON implementation for WebFormatter interface.
func (r ReturnRate) JSON() ([]byte, error) {
	return json.Marshal(r)
}

// XML implementation for WebFormatter interface.
func (r ReturnRate) XML() ([]byte, error) {
	return xml.Marshal(r)
}

func lookupRate(d Duration) uint {
	// TODO: implement lookupRate
	return uint(0)
}

// currentWeeklyRates is (the only) global variable containing rate information for given times
var currentWeeklyRates WeeklyRates

// WeeklyRates is a mapping of DailyRates for each day of the week
type WeeklyRates map[time.Weekday]DailyRates

// DailyRates is mapping between "minutes-since-midnight" and the associated rate that begins at that time
type DailyRates map[uint64]Rate

// Rate contains a price corresponding to a specific time-range and specific day of the week.
type Rate struct {
	Day         time.Weekday `json:"day"`
	StartMinute uint64       `json:"start"`
	EndMinute   uint64       `json:"end"`
	Price       uint         `json:"price"`
}

// NewWeeklyRates creates an empty WeeklyRates table that is ready to populate
func NewWeeklyRates() WeeklyRates {
	rates := make(WeeklyRates)
	rates[time.Monday] = make(DailyRates)
	rates[time.Tuesday] = make(DailyRates)
	rates[time.Wednesday] = make(DailyRates)
	rates[time.Thursday] = make(DailyRates)
	rates[time.Friday] = make(DailyRates)
	rates[time.Saturday] = make(DailyRates)
	rates[time.Sunday] = make(DailyRates)
	return rates
}

// Keys returns a sorted []uint64 array of start-time indexes (keys) for DailyRates.
//
// NOTE: When using Keys() with range, remember that it is actually the values
//       the of range that comprise the keys for the rates map.
//
//		Example:
//			for _, key := range rates.Keys() {
// 				rate = rates[key])
// 			}
//
func (d *DailyRates) Keys() []uint64 {
	// The reason for all this sort jankiness is because the sort package
	// does not natively support sorting of []uint64 arrays.
	//
	// As a workaround, sort as []int first, then convert back to []uint64
	var iKeys []int
	var uKeys []uint64
	for uKey := range *d {
		iKeys = append(iKeys, int(uKey))
	}

	// Make sorted []int
	sort.Ints(iKeys)

	// Make sorted []uint64
	for _, iKey := range iKeys {
		uKeys = append(uKeys, uint64(iKey))
	}

	return uKeys
}

// HasRate determines if a rate exists for the given time in units of minutes-since-midnight, returns
// true if a rate exists for the given time, false otherwise.
func (d *DailyRates) HasPrice(minuteSinceMidnight uint64) bool {
	m := minuteSinceMidnight
	if _, err := (*d).AtMinuteSinceMidnight(m); err != nil {
		return false
	} else {
		return true
	}
}

// AtMinuteSinceMidnight determines if a rate exists for the given time in units of minutes-since-midnight. Returns
// the price if a rate exists for the given time, otherwise an error is returned.
func (d *DailyRates) AtMinuteSinceMidnight(m uint64) (Rate, error) {
	for _, key := range (*d).Keys() {
		rate := (*d)[key]
		if m >= rate.StartMinute && m < rate.EndMinute {
			return rate, nil
		}
	}
	return Rate{}, fmt.Errorf("no rate exists for minute: %v", m)
}

// Update accepts a JSON byte array to update the weekly rates.
func (rates *WeeklyRates) Update(jsonNewRates []byte) error {
	// Unmarshalling into ConfigRates will catch any JSON formatting errors early on.
	config := ConfigRates{}
	err := json.Unmarshal(jsonNewRates, &config)
	if err != nil {
		fmt.Printf("could not parse JSON to update rates : %v", err.Error())
		return err
	}

	// Update weekly rates with each new rate.
	for _, newRateConfig := range config.Rates {
		if err = updateRate(newRateConfig, rates); err != nil {
			return fmt.Errorf("could not update rate: %v", err.Error())
		}
	}

	return nil
}

// LookupByDuration returns a price for given start and end timestamps in RFC3339 format.
// If the price is not available, and error is returned with 0 price.
func (weekRates *WeeklyRates) Lookup(start string, end string) (uint, error) {
	duration, err := ParseDuration(start, end)
	if err != nil {
		return 0, err
	}

	return weekRates.LookupByDuration(duration)
}

// updateRate is a helper function for the WeeklyRates.Update() method that attempts to
// update the rate configuration for a single time-slot.
func updateRate(rate ConfigRate, rates *WeeklyRates) error {
	start, end, err := TimeRangeFromConfigString(rate.Times)
	if err != nil {
		return err
	}

	days := strings.Split(rate.Days, ",")
	for _, day := range days {
		weekday, err := WeekdayFromString(day)
		if err != nil {
			return err
		}

		// Create new rate for time-window
		newRate := Rate{Day: weekday, StartMinute: start, EndMinute: end, Price: rate.Price}
		if err = rates.ConflictsWith(newRate); err != nil {
			return fmt.Errorf("new rate presents a conflict %v: %v", newRate, err.Error())
		}

		// Insert new rate for the appropriate weekday
		(*rates)[weekday][start] = newRate
	}

	return nil
}

// LookupByDuration returns a price for the time duration, if available.
func (weekRates *WeeklyRates) LookupByDuration(d Duration) (uint, error) {

	if d.Start.YearDay() != d.End.YearDay() {
		return 0, fmt.Errorf("start and end times must be on same day: start=%v, end=%v", d.Start.YearDay(), d.End.YearDay())
	}

	dayRates := (*weekRates)[d.Start.Weekday()]
	startMin := MinutesSinceMidnightFromTime(d.Start)
	endMin := MinutesSinceMidnightFromTime(d.End)

	startRate, err := dayRates.AtMinuteSinceMidnight(startMin)
	if err != nil {
		return 0, fmt.Errorf("rate unavailable: %v", err.Error())
	}

	endRate, err := dayRates.AtMinuteSinceMidnight(endMin)
	if err != nil {
		return 0, fmt.Errorf("rate unavailable: %v", err.Error())
	}

	// Require that rates are in the same range, even if the numeric price
	// is equal.
	if startRate == endRate {
		return startRate.Price, nil
	}

	return 0, fmt.Errorf("rate not in same time range")
}

// ConflictsWith determines if a new Rate will overlap with any existing Rate in WeeklyRates, and
// returns true if a conflict exists, false if no conflict.
func (weekRates *WeeklyRates) ConflictsWith(newRate Rate) error {
	dayRates := (*weekRates)[newRate.Day]

	rate, err := dayRates.AtMinuteSinceMidnight(newRate.StartMinute)
	if err == nil {
		return fmt.Errorf("rate already exists %v", rate)
	}
	rate, err = dayRates.AtMinuteSinceMidnight(newRate.EndMinute)
	if err == nil {
		return fmt.Errorf("rate already exists %v", rate)
	}
	return nil
}

// RateHandleFunc provides an endpoint to that echos back both a start and end timestamp
// in RFC3339 format, after parsing and computing duration.
//
// Example:
// 		curl  "http://localhost:8080/api/duration?start=2015-07-01T07%3A00%3A00Z&end=2015-07-01T12%3A00%3A00Z"
func RateHandleFunc(w http.ResponseWriter, r *http.Request) {
	InitializeResponse(&w, r) // Required before WriteResponse

	// Calculate duration from start to end
	duration, err := DurationFromHTTPRequest(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(Error{err.Error()}, &w)
		return
	}

	// Lookup the ReturnRate
	rate := ReturnRate{Start: duration.Start, End: duration.End, Price: lookupRate(duration)}

	// Respond with price
	w.WriteHeader(http.StatusOK)
	err = WriteResponse(rate, &w)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteResponse(Error{err.Error()}, &w)
		return
	}
}
