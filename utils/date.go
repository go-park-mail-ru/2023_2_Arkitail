package utils

import (
	"encoding/json"
	"strings"
	"time"
)

type Date time.Time

func (dateValue *Date) UnmarshalJSON(b []byte) error {
	cleanString := strings.Trim(string(b), "\"")
	time, err := time.Parse("2006-01-02", cleanString)
	if err != nil {
		return err
	}
	*dateValue = Date(time)
	return nil
}

func (dateValue Date) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(dateValue))
}

func (dateValue Date) Format(str string) string {
	time := time.Time(dateValue)
	return time.Format(str)
}
