package models

import (
	"fmt"
	"strings"
	"time"
)

type DataSource int

const (
    SourceNone DataSource = iota
    SourceMemory
    SourceRedis
)

type CustomTime struct {
	time.Time
}

type CustomDateTime struct {
	time.Time
}

const timeFormat = "15:04:05"
const isoLayout = "2006-01-02T15:04:05"

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

func (cdt *CustomDateTime) UnmarshalJSON(data []byte) error {
	str := strings.Trim(string(data), "\"")

	if str == "" || str == "0001-01-01T00:00:00" {
		cdt.Time = time.Time{} // Zero time
		return nil
	}

	t, err := time.Parse(isoLayout, str)
	if err != nil {
		return fmt.Errorf("CustomDateTime: format no reconegut: %s, error: %w", str, err)
	}

	cdt.Time = t
	return nil
}

func (cdt CustomDateTime) MarshalJSON() ([]byte, error) {
	if cdt.IsZero() {
		return []byte(`""`), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", cdt.Format(isoLayout))), nil
}