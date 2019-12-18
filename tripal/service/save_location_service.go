package service

import (
	. "../model"
	"context"
	. "github.com/common-go/search"
	. "github.com/common-go/service"
)

type SaveLocationService interface {
	GenericService
	SearchService
	GetLocationsOfUser(ctx context.Context, userId string) ([]Location, error)
}
