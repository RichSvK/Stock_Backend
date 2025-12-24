package route

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	router.Use(cors.Default())

	validate := validator.New()

	BalanceRoute(router, db, validate)
	StockWebRoute(router, db, validate)
	IpoRoute(router, db, validate)
	BrokerRoute(router, db, validate)
	StockRoute(router, db, validate)
	return router
}
