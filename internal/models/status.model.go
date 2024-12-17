package models

import (
	"time"

	"github.com/google/uuid"
)

type Status struct {
	Id               uuid.UUID `json:"Id"`
	Name             string    `json:"Name"`
	OperatorsAllowed bool      `json:"OperatorsAllowed"`
	Closed           bool      `json:"Closed"`
	Stopped          bool      `json:"Stopped"`
	Color            string    `json:"Color"`
	StartTime        time.Time `json:"StartTime"`
}

type ChangeStatusRequest struct {
	WorkcenterId uuid.UUID `json:"WorkcenterId"`
	StatusId uuid.UUID `json:"StatusId"`
}