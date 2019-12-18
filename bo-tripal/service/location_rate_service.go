package service

import (
	. "github.com/common-go/search"
	. "github.com/common-go/service"
)

type LocationRateService interface {
	ViewService
	SearchService
}
