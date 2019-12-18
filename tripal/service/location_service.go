package service

import (
	. "../model"
	"context"
	. "github.com/common-go/search"
	. "github.com/common-go/service"
)

type LocationService interface {
	GenericService
	SearchService
	GetByUrlId(ctx context.Context, urlId string) (*Location, error)
	RateLocation(ctx context.Context, objRate LocationRate) (bool, error)
	SaveLocation(ctx context.Context, userId string, locationId string) (bool, error)
	RemoveLocation(ctx context.Context, userId string, locationId string) (bool, error)
	GetLocationsOfUser(ctx context.Context, userId string) ([]Location, error)
	GetLocationByTypeInRadius(ctx context.Context, type1 string, raidus int) ([]Location, error)
}
