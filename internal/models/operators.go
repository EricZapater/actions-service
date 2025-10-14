package models

import "github.com/google/uuid"

// OperatorRequest representa una petició d'operador per clock in/out
type OperatorRequest struct {
	WorkcenterID uuid.UUID `json:"workcenterId" binding:"required" example:"123e4567-e89b-12d3-a456-426614174000"`
	OperatorID   uuid.UUID `json:"operatorId" binding:"required" example:"123e4567-e89b-12d3-a456-426614174001"`
}

// OperatorClockInDTO representa les dades d'un clock in d'operador
type OperatorClockInDTO struct {
	OperatorId   uuid.UUID `json:"operatorId" example:"123e4567-e89b-12d3-a456-426614174001"`
	WorkcenterId uuid.UUID `json:"workcenterId" example:"123e4567-e89b-12d3-a456-426614174000"`
	Timestamp    string    `json:"timestamp" example:"2025-10-14T10:30:00Z"`
}

// OperatorResponse representa la resposta amb les dades bàsiques d'un operador
type OperatorResponse struct {
	OperatorId     uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174001"`
	Code           string    `json:"code" example:"OP001"`
	Name           string    `json:"name" example:"Joan"`
	Surname        string    `json:"surname" example:"Garcia"`
	OperatorTypeID uuid.UUID `json:"operatorTypeId" example:"123e4567-e89b-12d3-a456-426614174002"`
}

// OperatorTypeResponse representa la resposta amb les dades d'un tipus d'operador
type OperatorTypeResponse struct {
	OperatorTypeId uuid.UUID `json:"id" example:"123e4567-e89b-12d3-a456-426614174002"`
	Name           string    `json:"name" example:"Operador Senior"`
	Description    string    `json:"description" example:"Operador amb més de 5 anys d'experiència"`
	Cost           float64   `json:"cost" example:"25.50"`
}

// OperatorDTO representa les dades completes d'un operador amb el seu tipus
type OperatorDTO struct {
	OperatorID              uuid.UUID `json:"operatorId" example:"123e4567-e89b-12d3-a456-426614174001"`
	OperatorCode            string    `json:"operatorCode" example:"OP001"`
	OperatorName            string    `json:"operatorName" example:"Joan"`
	OperatorSurname         string    `json:"operatorSurname" example:"Garcia"`
	OperatorTypeID          uuid.UUID `json:"operatorTypeId" example:"123e4567-e89b-12d3-a456-426614174002"`
	OperatorTypeName        string    `json:"operatorTypeName" example:"Operador Senior"`
	OperatorTypeDescription string    `json:"operatorTypeDescription" example:"Operador amb més de 5 anys d'experiència"`
	OperatorTypeCost        float64   `json:"operatorTypeCost" example:"25.50"`
}