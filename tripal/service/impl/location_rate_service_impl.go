package impl

import (
	"reflect"

	"../../model"
	m "github.com/common-go/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

type LocationRateServiceImpl struct {
	*m.DefaultSearchService
}

func NewLocationRateServiceImpl(db *mongo.Database, searchBuilder m.SearchResultBuilder) *LocationRateServiceImpl {
	var model model.LocationRate
	typeOfModel := reflect.TypeOf(model)
	r := m.NewDefaultSearchService(db, typeOfModel, "locationRate", searchBuilder)
	return &LocationRateServiceImpl{r}
}
