package models

type WorkOrderPhaseAndStatusRequest struct {
	WorkcenterID     string  `json:"WorkcenterId"`
	WorkOrderPhaseId string  `json:"WorkOrderPhaseId"`
	MachineStatusId  *string `json:"MachineStatusId"`
	TimeStamp        *string `json:"timestamp"`
}

type WorkOrderDTO struct {
	WorkOrderPhaseId          string `json:"workOrderPhaseId"`
	WorkOrderCode             string `json:"workOrderCode"`
	WorkOrderPhaseCode        string `json:"workOrderPhaseCode"`
	WorkOrderPhaseDescription string `json:"workOrderPhaseDescription"`
	PlannedQuantity           int    `json:"plannedQuantity"`
	ReferenceCode             string `json:"referenceCode"`
	ReferenceDescription      string `json:"referenceDescription"`
	StartTime                 string `json:"startTime"`
}

type WorkOrderPhaseResponse struct {
	WorkOrderCode             string `json:"workOrderCode"`
	WorkOrderPhaseCode        string `json:"workOrderPhaseCode"`
	WorkOrderPhaseDescription string `json:"workOrderPhaseDescription"`
	PlannedQuantity           int    `json:"plannedQuantity"`
	ReferenceCode             string `json:"referenceCode"`
	ReferenceDescription      string `json:"referenceDescription"`
}