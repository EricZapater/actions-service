package models

type WorkOrderPhaseAndStatusRequest struct {
	WorkcenterID     string  `json:"WorkcenterId"`
	WorkOrderPhaseId string  `json:"WorkOrderPhaseId"`
	MachineStatusId  *string `json:"MachineStatusId"`
	TimeStamp        *string `json:"timestamp"`
}

type WorkOrderDTO struct {
	WorkcenterID     string `json:"WorkcenterId"`
	WorkOrderPhaseId string `json:"WorkOrderPhaseId"`
	MachineStatusId  string `json:"MachineStatusId"`
	StartTime        string `json:"StartTime"`
}