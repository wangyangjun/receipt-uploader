package service

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"

	"github.com/wangyangjun/receipt-uploader/graph/model"
)

func CreateRecept(receipt model.Receipt) error {
	file, err := ioutil.ReadFile(receiptFilePath)
	if err != nil {
		log.Println(err)
		return errors.New("Can not open receipt file")
	}

	receiptInternal := model.ReceiptInternal{
		ID:          receipt.ID,
		ImageName:   receipt.ImageName,
		UserID:      receipt.User.ID,
		DateCreated: receipt.DateCreated,
	}

	data := []model.ReceiptInternal{}
	json.Unmarshal(file, &data)

	data = append(data, receiptInternal)

	// Preparing the data to be marshalled and written.
	dataBytes, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		return errors.New("Receipt serialization failed")
	}

	err = ioutil.WriteFile(receiptFilePath, dataBytes, 0644)
	if err != nil {
		log.Println(err)
		return errors.New("Save receipt failed")
	}
	return nil
}

func GetAllRecepts() ([]*model.Receipt, error) {
	file, err := ioutil.ReadFile(receiptFilePath)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Can not open receipt file")
	}
	data := []model.ReceiptInternal{}
	json.Unmarshal(file, &data)

	receipts := []*model.Receipt{}
	for i := range data {
		user, _ := GetUserById(data[i].UserID)
		receipt := model.Receipt{
			ID:          data[i].ID,
			ImageName:   data[i].ImageName,
			ImageURL:    "http://localhost:8080/" + "images/" + data[i].ImageName,
			User:        user,
			DateCreated: data[i].DateCreated,
		}
		receipts = append(receipts, &receipt)
	}

	return receipts, nil
}

func GetReceptByID(id string) (*model.Receipt, error) {
	receipts, err := GetAllRecepts()
	if err == nil {
		for i := range receipts {
			if id == (*receipts[i]).ID {
				return receipts[i], nil
			}
		}
	}
	return nil, errors.New("No receipt with such id")
}
