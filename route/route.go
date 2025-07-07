package route

import (
	"backend/controller"
	"backend/repository"
	"backend/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	BalanceRoute(router)
	StockWebRoute(router)
	IpoRoute(router)
	BrokerRoute(router)
	return router
}

func BalanceRoute(router *gin.Engine) {
	balanceRepository := repository.NewBalanceRepository()
	balanceService := service.NewBalanceService(balanceRepository)
	balanceController := controller.NewBalanceController(balanceService)
	router.GET("/export", balanceController.ExportBalanceController)
	router.GET("/balance/:code", balanceController.GetBalanceChart)
	router.POST("/balance/upload", balanceController.Upload)
	router.GET("/balance/scriptless", balanceController.GetScriptlessChange)
}

func StockWebRoute(router *gin.Engine) {
	stockWebRepository := repository.NewStockWebRepository()
	stockWebService := service.NewStockWebService(stockWebRepository)
	stockWebController := controller.NewStockWebController(stockWebService)
	router.GET("/links", stockWebController.GetLinks)
	router.GET("/links/:category", stockWebController.GetLinkReference)
}

func IpoRoute(router *gin.Engine) {
	ipoRepository := repository.NewIpoRepository()
	ipoService := service.NewIpoService(ipoRepository)
	ipoController := controller.NewIpoController(ipoService)

	router.GET("/ipo", ipoController.GetAllIpo)
	router.GET("/ipo/value", ipoController.GetIpoByValue)
	router.GET("/ipo/underwriter/:underwriter", ipoController.GetIpoByUnderwriter)
	router.POST("/ipo/condition", ipoController.GetIpoByCondition)
}

func BrokerRoute(router *gin.Engine) {
	brokerRepository := repository.NewBrokerRepository()
	brokerService := service.NewBrokerService(brokerRepository)
	brokerController := controller.NewBrokerController(brokerService)
	router.GET("/brokers", brokerController.GetBrokers)
}
