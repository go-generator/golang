package impl

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

type ActivityLogServiceImpl struct {
	database   *mongo.Database
	collection *mongo.Collection
}

func NewActivityLogServiceImpl(database *mongo.Database) *ActivityLogServiceImpl {
	return &ActivityLogServiceImpl{database, database.Collection("user_activity_log")}
}

func (s *ActivityLogServiceImpl) SaveLog(ctx context.Context, resource string, action string, success bool, desc string) error {

	return nil
}

func getFieldFromContext(field string, ctx context.Context) string {
	token := ctx.Value("authorization")
	if authorizationToken, ok := token.(map[string]interface{}); ok {
		return fmt.Sprint(authorizationToken[field])
	}
	return ""
}
