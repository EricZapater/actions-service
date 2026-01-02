package models

type WorkOrderPhaseAndStatusRequest struct {
	WorkcenterID     string  `json:"WorkcenterId"`
	WorkOrderPhaseId string  `json:"WorkOrderPhaseId"`
	MachineStatusId  *string `json:"MachineStatusId"`
	TimeStamp        *string `json:"timestamp"`
}

type WorkOrderDTO struct {
	WorkOrderPhaseId string `json:"WorkOrderPhaseId"`
	StartTime        string `json:"StartTime"`
}