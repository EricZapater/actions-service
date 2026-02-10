package validator

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

// Service defines the validator service interface
type Service interface {
	// ValidateStatusChange validates and executes actions when status changes
	ValidateStatusChange(ctx context.Context, workcenterID, newStatusID string) error
	
	// ValidateOperatorClockIn validates if an operator can clock in
	ValidateOperatorClockIn(ctx context.Context, operatorID, workcenterID string) error
	
	// ValidateWorkOrderPhaseIn validates if a work order phase can start
	ValidateWorkOrderPhaseIn(ctx context.Context, phaseID, workcenterID string) error
}

type service struct {
	operatorPort       OperatorPort
	statusRepo         StatusRepository
	workOrderPhasePort WorkOrderPhasePort
	workcenterPort     WorkcenterPort
}

// NewValidatorService creates a new validator service
func NewValidatorService(
	operatorPort OperatorPort,
	statusRepo StatusRepository,
	workOrderPhasePort WorkOrderPhasePort,
	workcenterPort WorkcenterPort,
) Service {
	return &service{
		operatorPort:       operatorPort,
		statusRepo:         statusRepo,
		workOrderPhasePort: workOrderPhasePort,
		workcenterPort:     workcenterPort,
	}
}

// ValidateStatusChange validates and executes necessary actions when status changes
func (s *service) ValidateStatusChange(ctx context.Context, workcenterID, newStatusID string) error {
	// 1. Get the new status
	key := fmt.Sprintf("%s",newStatusID)
	newStatus, _, err := s.statusRepo.FindByID(ctx, key)
	if err != nil {
		return fmt.Errorf("failed to get new status: %w", err)
	}

	// 2. Get current workcenter state
	wc, err := s.workcenterPort.GetWorkcenterDTO(ctx, workcenterID)
	if err != nil {
		return fmt.Errorf("failed to get workcenter: %w", err)
	}

	// 3. ⭐ VALIDATE AND EXECUTE: Check operators
	if !newStatus.OperatorsAllowed && len(wc.Operators) > 0 {
		log.Printf("⚠️  Status change: New status doesn't allow operators. Clocking out %d operators from workcenter %s\n",
			len(wc.Operators), workcenterID)

		for _, op := range wc.Operators {
			if err := s.operatorPort.ClockOut(ctx, op.OperatorID.String(), workcenterID); err != nil {
				log.Printf("❌ Error clocking out operator %s: %v", op.OperatorID, err)
				// Continue with other operators even if one fails
			}
		}
	}

	// TODO: Add WorkOrder validation when that functionality is implemented
	// if !newStatus.WorkOrderAllowed && len(wc.WorkOrders) > 0 { ... }

	return nil
}

// ValidateOperatorClockIn validates if an operator can clock in to a workcenter
func (s *service) ValidateOperatorClockIn(ctx context.Context, operatorID, workcenterID string) error {
	// 1. Get current workcenter state
	wc, err := s.workcenterPort.GetWorkcenterDTO(ctx, workcenterID)
	if err != nil {
		return fmt.Errorf("failed to get workcenter: %w", err)
	}

	// 2. ⭐ VALIDATE: Check if status allows operators
	if !wc.StatusOperatorsAllowed {
		return NewValidationError(
			http.StatusForbidden,
			fmt.Sprintf("operators not allowed in current status '%s'", wc.StatusName),
			nil,
		)
	}

	return nil
}

// ValidateWorkOrderPhaseIn validates if a work order phase can start on a workcenter
// TODO: Implement when WorkOrder functionality is added
func (s *service) ValidateWorkOrderPhaseIn(ctx context.Context, phaseID, workcenterID string) error {
	// Not implemented yet - WorkOrders not in main branch
	return nil
}
