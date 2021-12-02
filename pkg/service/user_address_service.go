package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/6156-DonaldDuck/composition/pkg/config"
	"github.com/6156-DonaldDuck/composition/pkg/model"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"sync"
)

func GetUserAddressById(userId uint) (model.UserAddress, error) {
	user := model.User{}
	address := model.Address{}

	var wg sync.WaitGroup
	wg.Add(2)

	go func(userId uint, user *model.User) {
		url := fmt.Sprintf("http://ec2-52-14-16-222.us-east-2.compute.amazonaws.com:8080/api/v1/users/%d", userId)
		resp, err := http.Get(url)
		if err != nil {
			log.Errorf("[service.GetUserAddressById] error occurred while getting user with id %v, err=%v\n", userId, err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &user)
		if err != nil {
			log.Errorf("[service.GetUserAddressById] error occurred while getting user with id %v, err=%v\n", userId, err)
		}
		wg.Done()
	}(userId, &user)

	go func(userId uint, address *model.Address) {
		url := fmt.Sprintf("http://ec2-52-14-16-222.us-east-2.compute.amazonaws.com:8085/api/v1/users/%d/address", userId)
		resp, err := http.Get(url)
		if err != nil {
			log.Errorf("[service.GetAddressByUserId] error occurred while getting address by user id %v, err=%v\n", userId, err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		json.Unmarshal(body, &address)
		if err != nil {
			log.Errorf("[service.GetAddressByUserId] error occurred while getting address by user id %v, err=%v\n", userId, err)
		}
		wg.Done()
	}(userId, &address)
	
	wg.Wait()

	useraddress := model.UserAddress{
		FirstName: user.FirstName,
		LastName: user.LastName,
		PhoneNumber: user.PhoneNumber,
		Email: user.Email,
		StreetName1: address.StreetName1,
		StreetName2: address.StreetName2,
		City: address.City,
		Region: address.Region,
		CountryCode: address.CountryCode,
		PostalCode: address.PostalCode,
	}

	return useraddress, nil
}

func GetUserById(userId string) (model.User, error) {
	user := model.User{}
	userUrl := config.Configuration.UserEndpoint+config.Configuration.BaseURL+"/users/"+userId
	resp, err := http.Get(userUrl)
	if err != nil {
		log.Errorf("[service.GetUserById] error occurred while getting user with id %v, err=%v\n", userId, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &user)
	if err != nil {
		log.Errorf("[service.GetUserById] error occurred while getting user with id %v, err=%v\n", userId, err)
	}

	return user, nil
}

func GetAddressByUserId(userId string) (model.Address, error) {
	address := model.Address{}
	addressUrl := config.Configuration.AddressEndpoint+config.Configuration.BaseURL+"/users/"+userId +"/address"
	resp, err := http.Get(addressUrl)
	if err != nil {
		log.Errorf("[service.GetAddressByUserId] error occurred while getting user with id %v, err=%v\n", userId, err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &address)
	if err != nil {
		log.Errorf("[service.GetAddressByUserId] error occurred while getting user with id %v, err=%v\n", userId, err)
	}
	return address, nil
}

func CreateUser(user model.User) (string, error) {
	body, _ := json.Marshal(user)
	userUrl := config.Configuration.UserEndpoint+config.Configuration.BaseURL+"/users"
	resp, err := http.Post(userUrl, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Errorf("[service.CreateUser] error occurred while getting user with err=%v\n", err)
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(bytes)
	userId := string(bytes)

	return userId, nil
}

func CreateAddress(address model.Address) (string, error) {
	body, _ := json.Marshal(address)
	addressUrl := config.Configuration.AddressEndpoint+config.Configuration.BaseURL+"/addresses"
	resp, err := http.Post(addressUrl, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Errorf("[service.CreateUser] error occurred while getting user with err=%v\n", err)
	}
	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(bytes)
	addressId := string(bytes)
	fmt.Println(addressId)

	return addressId, nil
}