package handler

import (
	"context"
	"github.com/google/uuid"
	"strings"
)

type ControlModelHandler struct {
	IdNames []string
}

func (c *ControlModelHandler) BuildToApprove(ctx context.Context, model interface{}) interface{} {
	panic("implement me")
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
	return obj
}

func (c *ControlModelHandler) BuildToUpdate(ctx context.Context, obj interface{}) interface{} {

	return obj
}

func (c *ControlModelHandler) BuildToPatch(ctx context.Context, obj interface{}) interface{} {
	return c.BuildToUpdate(ctx, obj)
}

func (c *ControlModelHandler) BuildToSave(ctx context.Context, obj interface{}) interface{} {
	return c.BuildToUpdate(ctx, obj)
}
