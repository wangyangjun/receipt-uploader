package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	"github.com/wangyangjun/receipt-uploader/graph/model"
)

func CreateUser(user model.User) error {
	file, err := ioutil.ReadFile(userFilePath)
	if err != nil {
		log.Println(err)
		return errors.New("Can not open users file")
	}

	data := []model.User{}
	json.Unmarshal(file, &data)

	data = append(data, user)

	// Preparing the data to be marshalled and written.
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return errors.New("User serialization failed")
	}

	err = ioutil.WriteFile(userFilePath, dataBytes, 0644)
	if err != nil {
		log.Println(err)
		return errors.New("Save user failed")
	}
	return nil
}

func GetAllUsers() ([]*model.User, error) {
	data := []model.User{}
	users := []*model.User{}

	file, err := ioutil.ReadFile(userFilePath)
	if err != nil {
		log.Println(err)
		return users, errors.New("Can not open users file")
	}

	json.Unmarshal(file, &data)
	for i := range data {
		users = append(users, &data[i])
	}

	return users, nil
}

func GetUserById(id string) (*model.User, error) {
	users, err := GetAllUsers()
	if err == nil {
		for i := range users {
			if id == (*users[i]).ID {
				return users[i], nil
			}
		}
	}
	return nil, errors.New("No user with such id")
}
