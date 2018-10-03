package context
import (
	. "../controller"
	. "../service"
	. "../repository"
	. "../config"
)
type AppContext struct {

}
var config = Config{}
var controller = MerchantController{}

var service = MerchantService{}

var dao = MerchantRepository{}
func (m* AppContext) GetMerchantController() MerchantController  {
	config.Read()
	dao.Database = config.Database
	dao.Server = config.Server
	dao.Connect()

	service.MerchantDao = &dao
	controller.MerchantService = &service
	return controller
}