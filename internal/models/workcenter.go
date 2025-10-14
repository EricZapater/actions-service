package models

import (
	"time"

	"github.com/google/uuid"
)

type WorkcenterDTO struct {
	WorkcenterID                uuid.UUID      `json:"workcenterId"`
	WorkcenterName              string         `json:"workcenterName"`
	WorkcenterDescription       string         `json:"workcenterDescription"`
	AreaID                      uuid.UUID      `json:"areaId"`
	AreaDescription             string         `json:"areaDescription"`
	ShiftID                     uuid.UUID      `json:"shiftId"`
	ShiftName                   string         `json:"shiftName"`
	ShiftDetailId               uuid.UUID      `json:"shiftDetailId"`
	ShiftDetailStartTime        CustomTime     `json:"shiftDetailStartTime"`
	ShiftDetailEndTime          CustomTime     `json:"shiftDetailEndTime"`
	ShiftDetailIsProductiveTime bool           `json:"shiftDetailsIsProductiveTime"`
	StatusID                    uuid.UUID      `json:"statusId"`
	StatusName                  string         `json:"statusName"`
	StatusOperatorsAllowed      bool           `json:"statusOperatorsAllowed"`
	StatusClosed                bool           `json:"statusClosed"`
	StatusStopped               bool           `json:"statusStopped"`
	StatusColor                 string         `json:"statusColor"`
	StatusStartTime             time.Time      `json:"statusStartTime"`
	Operators                   []OperatorDTO  `json:"operators"`
	//WorkOrders                  []WorkOrderDTO `json:"Workorders"`
}

type Workcenter struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AreaId      uuid.UUID `json:"areaId"`
	ShiftId     uuid.UUID `json:"shiftId"`
	Disabled    bool      `json:"disabled"`
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