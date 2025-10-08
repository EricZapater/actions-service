package models

import (
	"time"

	"github.com/google/uuid"
)

type WorkcenterDTO struct {
	WorkcenterID                uuid.UUID      `json:"WorkcenterId"`
	WorkcenterName              string         `json:"WorkcenterName"`
	WorkcenterDescription       string         `json:"WorkcenterDescription"`
	AreaID                      uuid.UUID      `json:"AreaId"`
	AreaDescription             string         `json:"AreaDescription"`
	ShiftID                     uuid.UUID      `json:"ShiftId"`
	ShiftName                   string         `json:"ShiftName"`
	ShiftDetailId               uuid.UUID      `json:"ShiftDetailId"`
	ShiftDetailStartTime        CustomTime     `json:"ShiftDetailStartTime"`
	ShiftDetailEndTime          CustomTime     `json:"ShiftDetailEndTime"`
	ShiftDetailIsProductiveTime bool           `json:"ShiftDetailsIsProductiveTime"`
	StatusID                    uuid.UUID      `json:"StatusId"`
	StatusName                  string         `json:"StatusName"`
	StatusOperatorsAllowed      bool           `json:"StatusOperatorsAllowed"`
	StatusClosed                bool           `json:"StatusClosed"`
	StatusStopped               bool           `json:"StatusStopped"`
	StatusColor                 string         `json:"StatusColor"`
	StatusStartTime             time.Time      `json:"StatusStartTime"`
	Operators                   []OperatorDTO  `json:"Operators"`
	//WorkOrders                  []WorkOrderDTO `json:"Workorders"`
}

type Workcenter struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AreaId      uuid.UUID `json:"areaId"`
	ShiftId     uuid.UUID `json:"shiftId"`
	Disabled    bool      `json:"Disabled"`
}

type CreateWorkcenterShiftDTO struct {
	WorkcenterID  uuid.UUID `json:"workcenterId"`
	ShiftDetailId uuid.UUID `json:"shiftDetailId"`
	StartTime     string    `json:"startTime"`
}

type Area struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	SiteId uuid.UUID `json:"siteId"`
}