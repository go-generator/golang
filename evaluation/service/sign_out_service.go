package service

import (
	"time"
)

type SignOutService interface {
	SignOut(token string, timeExpires time.Time) (bool, error)
}

