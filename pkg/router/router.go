package router

import (
	"errors"
	"github.com/6156-DonaldDuck/composition/pkg/config"
	"github.com/6156-DonaldDuck/composition/pkg/model"
	"github.com/6156-DonaldDuck/composition/pkg/router/middleware"
	"github.com/6156-DonaldDuck/composition/pkg/service"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func InitRouter() {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.GET("/api/v1/compositions/:userId", SyncGetUserAddressById)
	r.GET("/api/v1/user_address/:id", GetUserAddressById)
	r.POST("/api/v1/compositions", SyncPostUserAddressInfo)

	r.Run(":" + config.Configuration.Port)
}

func GetUserAddressById(c *gin.Context) {
	idStr := c.Param("id")
	userId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("[router.GetUserAddressById] failed to parse user id %v, err=%v\n", idStr, err)
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

func SyncGetUserAddressById(c *gin.Context) {
	composition :=model.Composition{}
	idStr := c.Param("userId")

	_, err := strconv.Atoi(idStr)
	if err != nil {
		log.Errorf("[router.GetComposedInfoById] failed to parse user id %v, err=%v\n", idStr, err)
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}

	// get user info
	user, err := service.GetUserById(idStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.Error(err)
		}
	}
	composition.User = user

	// get address info
	address, err := service.GetAddressByUserId(idStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, err.Error())
		} else {
			c.Error(err)
		}
	}
	composition.Address = address
	c.JSON(http.StatusOK, composition)
}

func SyncPostUserAddressInfo(c *gin.Context) {
	composition :=model.Composition{}
	if err := c.ShouldBind(&composition); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	userIdStr, err:= service.CreateUser(composition.User)

	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		log.Errorf("[router.SyncPostUserAddressInfo] failed to parse user id %v, err=%v\n", userIdStr, err)
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}

	composition.Address.UserId = uint(userId)
	addressIdStr, err:= service.CreateAddress(composition.Address)

	addressId, err := strconv.Atoi(addressIdStr)
	if err != nil {
		log.Errorf("[router.SyncPostUserAddressInfo] failed to parse address id %v, err=%v\n", addressIdStr, err)
		c.JSON(http.StatusBadRequest, "invalid user id")
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user": userId, "address": addressId})
}