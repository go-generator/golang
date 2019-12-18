package handler

import (
	"../model"
	"context"
	"github.com/common-go/echo"
	"github.com/google/uuid"
	"reflect"
	"strings"
	"time"
)

type ControlModelHandler struct {
	IdNames []string
}

func NewControlModelHandler(idNames []string) *ControlModelHandler {
	return &ControlModelHandler{idNames}
}

func getUserNameFromContext(ctx context.Context) string {
	token := ctx.Value("authorization")
	if authorizationToken, ok := token.(map[string]interface{}); ok {
		userName, _ := authorizationToken["userName"].(string)
		return userName
	}
	return ""
}

func generateId() string {
	id := uuid.New()
	return replaceAll(id.String(), "-", "")
}

func replaceAll(value string, strFind string, strReplace string) string {
	return strings.Replace(value, strFind, strReplace, -1)
}

func (c *ControlModelHandler) CheckSecurity(ctx context.Context, obj interface{}) bool {
	return true
}

func (c *ControlModelHandler) BuildToInsert(ctx context.Context, obj interface{}) interface{} {
	value := reflect.Indirect(reflect.ValueOf(obj))
	numField := value.NumField()

	if len(c.IdNames) > 0 {
		for _, field := range c.IdNames {
			if value.FieldByName(field).Len() != 0 {
				continue
			}
			server.SetField(obj, field, generateId())
		}
	}

	userName := getUserNameFromContext(ctx)
	for i := 0; i < numField; i++ {
		if controlModel, ok := value.Field(i).Interface().(*model.ControlModel); ok {
			if controlModel == nil {
				controlModel = &model.ControlModel{}
				value.Field(i).Set(reflect.ValueOf(controlModel))
			}
			controlModel.CtrlStatus = model.ControlStatusPending
			controlModel.ActionDate = time.Now()
			controlModel.ActedBy = userName
			controlModel.ActionStatus = model.ActionStatusCreated
			break
		}
	}
	return obj
}

func (c *ControlModelHandler) BuildToUpdate(ctx context.Context, obj interface{}) interface{} {
	value := reflect.Indirect(reflect.ValueOf(obj))
	numField := value.NumField()
	userName := getUserNameFromContext(ctx)
	for i := 0; i < numField; i++ {
		if controlModel, ok := value.Field(i).Interface().(*model.ControlModel); ok {
			if controlModel == nil {
				controlModel = &model.ControlModel{}
				value.Field(i).Set(reflect.ValueOf(controlModel))
			}
			controlModel.CtrlStatus = model.ControlStatusPending
			controlModel.ActionDate = time.Now()
			controlModel.ActedBy = userName
			controlModel.ActionStatus = model.ActionStatusUpdated
			break
		}
	}
	return obj
}

func (c *ControlModelHandler) BuildToPatch(ctx context.Context, obj interface{}) interface{} {
	return c.BuildToUpdate(ctx, obj)
}

func (c *ControlModelHandler) BuildToSave(ctx context.Context, obj interface{}) interface{} {
	return c.BuildToUpdate(ctx, obj)
}

func (c *ControlModelHandler) BuildToApprove(ctx context.Context, obj interface{}) interface{} {
	value := reflect.Indirect(reflect.ValueOf(obj))
	numField := value.NumField()
	userName := getUserNameFromContext(ctx)
	for i := 0; i < numField; i++ {
		if controlModel, ok := value.Field(i).Interface().(*model.ControlModel); ok {
			if controlModel == nil {
				controlModel = &model.ControlModel{}
				value.Field(i).Set(reflect.ValueOf(controlModel))
			}
			controlModel.ActionDate = time.Now()
			controlModel.ActedBy = userName
			controlModel.ActionStatus = model.ActionStatusUpdated
			break
		}
	}
	return obj
}
