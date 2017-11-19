package util

import (
	"strconv"
	"strings"
	"time"
)

// JSONTime - The `time.Time` type is not Unmarshaled from JSON given an unix epoch, by default.
type JSONTime time.Time

// MarshalJSON to marshal JSONTime to JSON
func (t JSONTime) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(time.Time(t).Unix(), 10)), nil
}

//UnmarshalJSON - marshal JSON value to JSONTime
func (t *JSONTime) UnmarshalJSON(s []byte) (err error) {
	r := strings.Replace(string(s), `"`, ``, -1)
	result, err := strconv.ParseInt(string(r), 10, 64)
	if err != nil {
		return err
	}

	// convert the unix epoch to a Time object
	*t = JSONTime(time.Unix(result/1000, 0))
	return nil
}

// String returns the string value of the time
func (t JSONTime) String() string {
	return time.Time(t).String()
}

//UTCString UTC format
func (t JSONTime) UTCString() string {
	return time.Time(t).UTC().String()
}
