package router

import (
	"github.com/6156-DonaldDuck/composition/pkg/config"
	"github.com/6156-DonaldDuck/composition/pkg/router/middleware"
	"github.com/gin-gonic/gin"
)

func InitRouter() {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.GET("/api/v1/user_address/:id", SyncGetUserAddressById)
	r.POST("/api/v1/user_address/", AsyncPostUserAddressInfo)

	r.Run(":" + config.Configuration.Port)
}

func SyncGetUserAddressById(c *gin.Context) {
	// TODO: synchronously get user and address info by id
}

func AsyncPostUserAddressInfo(c *gin.Context) {
	// TODO: asynchronously post user and address info
}
