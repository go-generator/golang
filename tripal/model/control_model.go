package model

import (
	"errors"
	"reflect"
	"time"
)

type ControlStatus string

const (
	ControlStatusPending  ControlStatus = "P"
	ControlStatusApproved ControlStatus = "A"
	ControlStatusReject   ControlStatus = "R"
)

type ActionStatus string

const (
	ActionStatusCreated ActionStatus = "C"
	ActionStatusUpdated ActionStatus = "U"
)

type ControlModel struct {
	ActedBy      string        `bson:"actedBy" json:"actedBy,omitempty" gorm:"type:varchar(50);column:acted_by"`
	ActionStatus ActionStatus  `bson:"actionStatus" json:"actionStatus,omitempty" gorm:"type:char(1);column:action_status"`
	CtrlStatus   ControlStatus `bson:"ctrlStatus" json:"ctrlStatus,omitempty" gorm:"type:char(1);column:ctrl_status"`
	ActionDate   time.Time     `bson:"actionDate" json:"actionDate,omitempty" gorm:"column:action_date"`
}

func IsExtendedFromControlModel(modelType reflect.Type) (error, bool) {
	var controlModel = reflect.New(modelType).Interface()
	value := reflect.Indirect(reflect.ValueOf(controlModel))
	numField := value.NumField()
	for i := 0; i < numField; i++ {
		if _, ok := value.Field(i).Interface().(*ControlModel); ok {
			return nil, true
		}
	}
	return errors.New(modelType.Name() + " isn't extended from *ControlModel struct!"), false
}
