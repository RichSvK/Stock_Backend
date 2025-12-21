package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	BalanceRoute(router, db)
	StockWebRoute(router, db)
	IpoRoute(router, db)
	BrokerRoute(router, db)
	StockRoute(router, db)
	return router
}
