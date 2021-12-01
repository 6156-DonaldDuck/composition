package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/6156-DonaldDuck/compositions/pkg/config"
	"github.com/6156-DonaldDuck/compositions/pkg/model"
	"github.com/6156-DonaldDuck/compositions/pkg/router/middleware"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"strconv"
)

func InitRouter() {
	r := gin.Default()
	r.Use(middleware.CORSMiddleware())

	r.GET("/api/v1/compositions/:userId", SyncGetUserAddressById)
	r.POST("/api/v1/compositions", AsyncPostUserAddressInfo)

	r.Run(":" + config.Configuration.Port)
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
	userUrl := config.Configuration.UserEndpoint+config.Configuration.BaseURL+"/users/"+idStr
	user := make(chan *http.Response)
	go SendGetAsync(userUrl, user)
	userResponse := <- user
	defer userResponse.Body.Close()
	bytes, _ := ioutil.ReadAll(userResponse.Body)
	jsonStr := string(bytes)
	err = json.Unmarshal([]byte(jsonStr), &composition.User)

	if err != nil {
		log.Errorf("[router.GetComposedInfoById] failed to parse user with id =%v, err=%v\n", idStr, err)
		c.JSON(http.StatusBadRequest, "invalid user")
		return
	}

	// get address info
	addressUrl := config.Configuration.AddressEndpoint+config.Configuration.BaseURL+"/addresses"+idStr+"/address"
	address := make(chan *http.Response)
	go SendGetAsync(addressUrl, address)
	addressResponse := <- address
	defer addressResponse.Body.Close()
	bytes, _ = ioutil.ReadAll(addressResponse.Body)
	jsonStr = string(bytes)
	err = json.Unmarshal([]byte(jsonStr), &composition.Address)

	// Temporarily comment this part. Because even when address is null, function should return a composition.

	//if err != nil {
	//	log.Errorf("[router.GetComposedInfoById] failed to parse address with its id =%v, err=%v\n", idStr, err)
	//	c.JSON(http.StatusBadRequest, "invalid address")
	//	return
	//}

	c.JSON(http.StatusOK, composition)
}

func AsyncPostUserAddressInfo(c *gin.Context) {
	composition :=model.Composition{}
	if err := c.ShouldBind(&composition); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}

	body, _ := json.Marshal(composition.User)
	userChan := make(chan *http.Response)
	// UserService
	userUrl := config.Configuration.UserEndpoint+config.Configuration.BaseURL+"/users"
	go SendPostAsync(userUrl, body, userChan)
	userResponse := <-userChan
	defer userResponse.Body.Close()
	bytes, err := ioutil.ReadAll(userResponse.Body)
	idStr := string(bytes)
	userId, err := strconv.Atoi(idStr)

	if err != nil {
		c.Error(err)
	} else {
		composition.Address.UserId = uint(userId)
		body, _ = json.Marshal(composition.Address)
		addressChan := make(chan *http.Response)
		// AddressService
		addressUrl := config.Configuration.AddressEndpoint+config.Configuration.BaseURL+"/addresses"
		go SendPostAsync(addressUrl, body, addressChan)
		addressResponse := <-addressChan
		defer addressResponse.Body.Close()
		bytes, _ := ioutil.ReadAll(addressResponse.Body)
		addressId := string(bytes)
		fmt.Println(addressId)

		c.JSON(http.StatusCreated, gin.H{"user": idStr, "address": addressId})
	}
}

func SendPostAsync(url string, body []byte, rc chan *http.Response) {
	response, err := http.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	rc <- response
}

func SendGetAsync(url string, rc chan *http.Response) {
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	rc <- response
}