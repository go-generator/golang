package service

import (
	. "../repository"
	. "../models"
	"github.com/reactivex/rxgo/observable"
)
type MerchantService struct {
	MerchantDao   *MerchantRepository
}
func (m* MerchantService) GetAll() ([]Merchant , error)  {
	return m.MerchantDao.FindAll()
}

func (m* MerchantService) GetAllObserve() observable.Observable  {
	return m.MerchantDao.FindAllObserver()
}
