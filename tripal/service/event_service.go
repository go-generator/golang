package service

import (
	. "../model"
	"context"
	. "github.com/common-go/search"
	. "github.com/common-go/service"
)

type EventService interface {
	GenericService
	SearchService
	GetEventByLocation(ctx context.Context, locationId string) ([]Event, error)
	GetEventByDate(ctx context.Context, date string) ([]Event, error)
}
