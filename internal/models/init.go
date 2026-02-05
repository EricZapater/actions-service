package models

type WorkcenterShiftDTO struct {
	WorkcenterID  string                     `json:"workcenterId"`
	ShiftDetailID string                     `json:"shiftDetailId"`
	StartTime     string                     `json:"startTime"`
	Details       []WorkcenterShiftDetailDTO `json:"details"`
}

type WorkcenterShiftDetailDTO struct {
	WorkCenterShiftId     string `json:"workcenterShiftId"`
	MachineStatusId       string `json:"machineStatusId"`
	MachineStatusReasonId string `json:"machineStatusReasonId"`
	OperatorId            string `json:"operatorId"`
	WorkOrderPhaseId      string `json:"workorderPhaseId"`
	StartTime             string `json:"startTime"`
	QuantityOk            int    `json:"quantityOk"`
	QuantityKo            int    `json:"quantityKo"`
}