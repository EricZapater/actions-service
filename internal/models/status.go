package models

import "github.com/google/uuid"

type StatusInRequest struct {
	StatusID     uuid.UUID `json:"machineStatusId"`
	WorkcenterID uuid.UUID `json:"workcenterId"`
	StatusReasonId *uuid.UUID `json:"machineStatusReasonId"`
	Timestamp	string    `json:"timestamp"`
}

type StatusDTORequest struct {
	StatusID     uuid.UUID `json:"statusId"`
	WorkcenterID uuid.UUID `json:"workcenterId"`
	StatusReasonId *uuid.UUID `json:"statusReasonId"`
}

type StatusResponse struct {
	StatusId uuid.UUID `json:"id"`
	Description string `json:"description"`
	Color string `json:"color"`
	Stopped bool `json:"stopped"`
	OperatorsAllowed bool `json:"operatorsAllowed"`
	IsDefault bool `json:"default"`
	Closed bool `json:"closed"`
}

type StatusCostResponse struct {
	StatusId     uuid.UUID `json:"machineStatusId"`
	WorkcenterId uuid.UUID `json:"workcenterId"`
	Cost float32 `json:"cost"`
}

type StatusDTO struct {	
	StatusId uuid.UUID `json:"statusId"`
	StatusReasonId *uuid.UUID `json:"statusReasonId"`
	Description string `json:"description"`
	Color string `json:"color"`
	Stopped bool `json:"stopped"`
	OperatorsAllowed bool `json:"operatorsAllowed"`	
	WorkOrdersAllowed bool `json:"workOrdersAllowed"`
	IsDefault bool `json:"default"`
	Closed bool `json:"closed"`
	Cost float32 `json:"cost"`
	StatusStartTime string `json:"statusStartTime"`
}