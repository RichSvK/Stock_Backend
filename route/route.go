package route

import (
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
	StockRoute(router)
	return router
}
