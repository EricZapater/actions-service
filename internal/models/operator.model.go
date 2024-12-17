package models

import "github.com/google/uuid"

type Operator struct {
}

type OperatorRequest struct {
	WorkcenterId uuid.UUID `json:"WorkcenterId"`
	OperatorId   uuid.UUID `json:"OperatorId"`
}