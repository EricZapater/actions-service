package models

import "github.com/google/uuid"

type OperatorRequest struct {
	WorkcenterId uuid.UUID `json:"WorkcenterId"`
	OperatorId   uuid.UUID `json:"OperatorId"`
}

type OperatorClockInDTO struct {
	OperatorId   uuid.UUID `json:"OperatorId"`
	WorkcenterId uuid.UUID `json:"WorkcenterId"`
	Timestamp   string    `json:"Timestamp"`
}

type OperatorResponse struct {
	OperatorId   uuid.UUID `json:"Id"`
	Code 	  string    `json:"Code"`
	Name 	  string    `json:"Name"`
	Surname    string    `json:"Surname"`
	OperatorTypeID uuid.UUID `json:"OperatorTypeId"`
}

type OperatorTypeResponse struct {
	OperatorTypeId uuid.UUID `json:"Id"`
	Name		  string    `json:"Name"`
	Description   string    `json:"Description"`
	Cost		 float64   `json:"Cost"`
}

type OperatorDTO struct {
	OperatorId   uuid.UUID `json:"OperatorId"`
	OperatorCode string    `json:"OperatorCode"`
	OperatorName string    `json:"OperatorName"`
	OperatorSurname string    `json:"OperatorSurname"`
	OperatorTypeId uuid.UUID `json:"OperatorTypeId"`
	OperatorTypeName string `json:"OperatorTypeName"`
	OperatorTypeDescription string `json:"OperatorTypeDescription"`
	OperatorTypeCost float64 `json:"OperatorTypeCost"`
}