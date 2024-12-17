package models

import (
	"fmt"
	"strings"
	"time"
)

type CustomTime struct {
	time.Time
}

const timeFormat = "15:04:05"

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	// Traiem les cometes del string JSON
	strTime := strings.Trim(string(data), "\"")
	// Parsegem el temps segons el format hh:mm:ss
	parsedTime, err := time.Parse(timeFormat, strTime)
	if err != nil {
		return fmt.Errorf("error parsing time: %w", err)
	}
	// Assignem el temps al CustomTime
	ct.Time = parsedTime
	return nil
}

// Implementem MarshalJSON per serialitzar CustomTime
func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", ct.Format(timeFormat))), nil
}