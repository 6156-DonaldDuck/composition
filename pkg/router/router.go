package router

import (
	"errors"
	"net/http"
	"strconv"
	"github.com/6156-DonaldDuck/composition/pkg/config"
	"github.com/6156-DonaldDuck/composition/pkg/service"
	"github.com/6156-DonaldDuck/composition/pkg/router/middleware"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func InitRouter() {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.GET("/api/v1/user_address/:id", SyncGetUserAddressById)
	r.POST("/api/v1/user_address/", AsyncPostUserAddressInfo)

	r.Run(":" + config.Configuration.Port)
}

func SyncGetUserAddressById(c *gin.Context) {
	idStr := c.Param("id")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("[router.SyncGetUserAddressById] failed to parse user id %v, err=%v\n", idStr, err)
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}
	user_address, err := service.GetUserAddressById(uint(userId))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.Error(err)
		}
	} else {
		c.JSON(http.StatusOK, user_address)
	}
}

func AsyncPostUserAddressInfo(c *gin.Context) {
	// TODO: asynchronously post user and address info
}
