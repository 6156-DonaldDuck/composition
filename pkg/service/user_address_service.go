package service

import (
	"github.com/6156-DonaldDuck/composition/pkg/model"
	log "github.com/sirupsen/logrus"
	"sync"
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
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