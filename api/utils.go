package api

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// WeekdayFromString returns a time.Weekday for the corresponding abbreviated day.
// Valid values for the day string include "mon", "tues", "wed", "thurs", "fri", "sat", and "sun".
func WeekdayFromString(day string) (time.Weekday, error) {
	var weekday time.Weekday

	switch day {
	case "mon":
		weekday = time.Monday
	case "tues":
		weekday = time.Tuesday
	case "wed":
		weekday = time.Wednesday
	case "thurs":
		weekday = time.Thursday
	case "fri":
		weekday = time.Friday
	case "sat":
		weekday = time.Saturday
	case "sun":
		weekday = time.Sunday
	default:
		return 0, fmt.Errorf("'%s' is not a recognized weekday", day)
	}

	return weekday, nil
}

func MinutesSinceMidnightFromTime(t time.Time) uint64 {
	return uint64(t.Minute() + 60*t.Hour())
}

// MinutesSinceMidnightFromString parses a time string in the form "0600" and returns total minutes since midnight.
// In the example of "0600" the value 3600 will be returned, since "0600" is 6*60 minutes since midnight.
func MinutesSinceMidnightFromString(s string) (uint64, error) {
	expectedLength := 4
	if len(s) != expectedLength {
		return 0, fmt.Errorf("start time must be %d characters: %s", expectedLength, s)
	}

	hourString := s[0:2]
	minuteString := s[2:4]

	hours, err := strconv.ParseUint(hourString, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid hours in time: %s", s)
	}

	minutes, err := strconv.ParseUint(minuteString, 10, 64)
	if err != nil || minutes >= 60 {
		return 0, fmt.Errorf("invalid minutes in time: %s", s)
	}

	totalMinutes := hours*60 + minutes
	return totalMinutes, nil
}

func TimeRangeFromConfigString(s string) (uint64, uint64, error) {
	times := strings.Split(s, "-")
	if len(times) != 2 {
		return 0, 0, fmt.Errorf("invalid time range format: %s", s)
	}

	start, err := MinutesSinceMidnightFromString(times[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid start time: %v", err.Error())
	}

	end, err := MinutesSinceMidnightFromString(times[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid end time: %s", err.Error())
	}

	return start, end, nil
}

func JSONFromRequestBody(r *http.Request) ([]byte, error) {
	// Require JSON Content-Type
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "json") {
		return nil, fmt.Errorf("invalid content type \"%v\" in request, \"Content-Type:application/json\" is required", contentType)
	}

	// Convert JSON body of request into []byte for unmarshalling
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return nil, err
	}

	return body, nil
}
