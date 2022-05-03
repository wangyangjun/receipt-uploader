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

	userMap := make(map[string]model.User)
	json.Unmarshal(file, &userMap)
	userMap[user.ID] = user
	userMapBytes, err := json.Marshal(userMap)
	if err != nil {
		log.Println(err)
		return errors.New("User data serialization failed")
	}
	err = ioutil.WriteFile(userFilePath, userMapBytes, 0644)

	if err != nil {
		log.Println(err)
		return errors.New("Save user failed")
	}
	return nil
}

func readUserMap() (map[string]model.User, error) {
	file, err := ioutil.ReadFile(userFilePath)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Can not open user file")
	}
	userMap := make(map[string]model.User)
	json.Unmarshal(file, &userMap)
	return userMap, nil
}

func GetAllUsers() ([]*model.User, error) {
	uerMap, err := readUserMap()
	if err != nil {
		return nil, err
	}
	users := []*model.User{}

	for _, user := range uerMap {
		users = append(users, &user)
	}

	return users, nil
}

func GetUserById(id string) (*model.User, error) {
	uerMap, err := readUserMap()
	if err != nil {
		return nil, err
	}
	if user, ok := uerMap[id]; ok {
		return &user, nil
	}
	return nil, errors.New("No user with such id")
}
