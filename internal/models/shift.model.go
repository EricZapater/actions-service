package models

import (
	"github.com/google/uuid"
)

type Shift struct {
	Id uuid.UUID `json:"id"`
	Name string `json:"name"`	
	ShiftDetail []ShiftDetail `json:"shiftDetail"`
}

type ShiftDetail struct {	
	Id uuid.UUID `json:"id"`
	StartTime CustomTime `json:"startTime"`
	EndTime CustomTime `json:"endTime"`
	IsProductiveTime bool `json:"isProductiveTime"`
}