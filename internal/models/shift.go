package models

import (
	"github.com/google/uuid"
)

type ShiftDTO struct {
	ID           uuid.UUID        `json:"id"`
	Name         string           `json:"name"`
	ShiftDetails []ShiftDetailDTO `json:"shift_details"`
}

type ShiftDetailDTO struct {	
	ID uuid.UUID `json:"id"`
	StartTime CustomTime `json:"startTime"`
	EndTime CustomTime `json:"endTime"`
	IsProductiveTime bool `json:"isProductiveTime"`
}