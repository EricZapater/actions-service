package bootstrap

import (
	"actions-service/internal/clients"
	"actions-service/internal/models"
	"actions-service/internal/operator"
	"actions-service/internal/shift"
	"actions-service/internal/status"
	"actions-service/internal/workorderphase"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Service struct {
	RedisRepo *RedisRepo
	client    clients.HttpBackendClient	
	statusService status.Service
	operatorService operator.Service
	shiftService shift.Service
	workorderphaseService workorderphase.Service
}

type service interface {
	InitDTO(ctx context.Context) error
}

func NewService(redisRepo *RedisRepo, client clients.HttpBackendClient, statusService status.Service, operatorService operator.Service, shiftService shift.Service, workorderphaseService workorderphase.Service) *Service {
	return &Service{
		RedisRepo: redisRepo,
		client:    client,
		statusService: statusService,
		operatorService: operatorService,
		shiftService: shiftService,
		workorderphaseService: workorderphaseService,
	}
}

func (s *Service) InitDTO(ctx context.Context) error {
	url := "/api/WorkcenterShift/Currents"
	response, err := s.client.DoGetRequest(ctx, url)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode > 299 {		
		return fmt.Errorf("failed to get currents: %s", response.Status)
	}
	var workcenterShifts []models.WorkcenterShiftDTO
	err = json.NewDecoder(response.Body).Decode(&workcenterShifts)
	if err != nil {
		return err
	}
	fmt.Println("WorkcenterShifts: ", len(workcenterShifts))
	var workcenters []models.WorkcenterDTO
	for _, wc := range workcenterShifts {
		fmt.Println("Workcenter: ", wc)
		workcenter, err := s.PopulateDTO(ctx, wc)
		if err != nil {
			log.Println(err)
			continue
		}
		workcenters = append(workcenters, workcenter)
	}
	err = s.RedisRepo.SetMultiple(ctx, workcenters)
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) PopulateDTO(ctx context.Context, wcs models.WorkcenterShiftDTO) (models.WorkcenterDTO, error) {
	//recuperar shift
	fmt.Println("WorkcenterShiftId: ", wcs.Details[0].WorkCenterShiftId)
	shift, err := s.shiftService.FindByID(ctx, wcs.Details[0].WorkCenterShiftId)
	if err != nil {
		return models.WorkcenterDTO{}, err
	}	
	shiftDetail, err := s.shiftService.FindShiftDetailByID(ctx, shift.ID.String(), wcs.ShiftDetailID)
	if err != nil {
		return models.WorkcenterDTO{}, err
	}
	workcenter := models.WorkcenterDTO{
		WorkcenterID: uuid.MustParse(wcs.WorkcenterID),
		ShiftID: shift.ID,
		ShiftName: shift.Name,
		ShiftDetailId: shiftDetail.ID,
		ShiftDetailStartTime: shiftDetail.StartTime,
		ShiftDetailIsProductiveTime: shiftDetail.IsProductiveTime,
	}
	for _, detail := range wcs.Details {
		status, err := s.statusService.FindByID(ctx, detail.MachineStatusId)
		if err != nil {
			return models.WorkcenterDTO{}, err
		}
		operator, err := s.operatorService.FindByID(ctx, detail.OperatorId)
		if err != nil {
			return models.WorkcenterDTO{}, err
		}
		workorderphase, err := s.workorderphaseService.FindByID(ctx, detail.WorkOrderPhaseId)
		if err != nil {
			return models.WorkcenterDTO{}, err
		}
		workorder := models.WorkOrderDTO{
			WorkOrderPhaseId: detail.WorkOrderPhaseId,			
			PlannedQuantity: workorderphase.PlannedQuantity,
			WorkOrderCode: workorderphase.WorkOrderCode,
			WorkOrderPhaseCode: workorderphase.WorkOrderPhaseCode,
			WorkOrderPhaseDescription: workorderphase.WorkOrderPhaseDescription,
			ReferenceCode: workorderphase.ReferenceCode,
			ReferenceDescription: workorderphase.ReferenceDescription,
			QuantityOk: detail.QuantityOk,
			QuantityKo: detail.QuantityKo,
		}

		var statusReasonId uuid.UUID
		if detail.MachineStatusReasonId != ""  {
			statusReasonId = uuid.MustParse(detail.MachineStatusReasonId)
		}
		layout := "2006-01-02 15:04:05"
		workcenter.StatusID = status.StatusId
		workcenter.StatusName = status.Description
		workcenter.StatusReasonId = &statusReasonId
		workcenter.StatusOperatorsAllowed = status.OperatorsAllowed
		workcenter.StatusWorkOrdersAllowed = status.WorkOrdersAllowed
		workcenter.StatusClosed = status.Closed
		workcenter.StatusStopped = status.Stopped
		workcenter.StatusColor = status.Color
		workcenter.StatusStartTime, _ = time.Parse(layout, detail.StartTime)
		workcenter.Operators = append(workcenter.Operators, operator)
		workcenter.WorkOrders = append(workcenter.WorkOrders, workorder)
	}
	return workcenter, nil
}