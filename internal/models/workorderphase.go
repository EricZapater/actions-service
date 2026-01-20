package models

type WorkOrderPhaseAndStatusRequest struct {
	WorkcenterID     string  `json:"workcenterId"`
	WorkOrderPhaseId string  `json:"workOrderPhaseId"`
	MachineStatusId  *string `json:"machineStatusId"`
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
	QuantityOk                int    `json:"quantityOk"`
	QuantityKo                int    `json:"quantityKo"`
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

type WorkOrderPhaseOutRequest struct {
	WorkcenterID         string  `json:"workcenterId"`
	WorkOrderPhaseId     string  `json:"workOrderPhaseId"`
	NextWorkOrderPhaseId *string `json:"nextWorkOrderPhaseId"`
	WorkOrderStatusId    *string `json:"workOrderStatusId"`
	NextMachineStatusId  *string `json:"nextMachineStatusId"`
	QuantityOk           *int    `json:"quantityOk"`
	QuantityKo           *int    `json:"quantityKo"`
	TimeStamp            *string `json:"timestamp"`
}

// BackendWorkOrderPhaseOutRequest is the request model for the backend API WorkOrderPhase/Out endpoint
type BackendWorkOrderPhaseOutRequest struct {
	WorkcenterID         string  `json:"workcenterId"`
	WorkOrderPhaseId     string  `json:"workOrderPhaseId"`
	TimeStamp            string  `json:"timestamp"`
	WorkOrderStatusId    *string `json:"workOrderStatusId,omitempty"`
	NextWorkOrderPhaseId *string `json:"nextWorkOrderPhaseId,omitempty"`
	NextMachineStatusId  *string `json:"nextMachineStatusId,omitempty"`
}

type WorkOrderPhaseQuantitiesRequest struct {
	WorkcenterID     string `json:"workcenterId"`
	WorkOrderPhaseId string `json:"workOrderPhaseId"`
	QuantityOk       int    `json:"quantityOk"`
	QuantityKo       int    `json:"quantityKo"`
}