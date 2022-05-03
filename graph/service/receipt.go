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

	receiptMap := make(map[string]model.ReceiptInternal)
	json.Unmarshal(file, &receiptMap)
	receiptMap[receiptInternal.ID] = receiptInternal
	receiptMapBytes, err := json.Marshal(receiptMap)
	if err != nil {
		log.Println(err)
		return errors.New("Receipt serialization failed")
	}
	err = ioutil.WriteFile(receiptFilePath, receiptMapBytes, 0644)
	if err != nil {
		log.Println(err)
		return errors.New("Save receipt failed")
	}
	return nil
}

func readReceiptMap() (map[string]model.ReceiptInternal, error) {
	file, err := ioutil.ReadFile(receiptFilePath)
	if err != nil {
		log.Println(err)
		return nil, errors.New("Can not open receipt file")
	}
	receiptMap := make(map[string]model.ReceiptInternal)
	json.Unmarshal(file, &receiptMap)
	return receiptMap, nil
}

func GetAllReceipts() ([]*model.Receipt, error) {
	receiptMap, err := readReceiptMap()
	if err != nil {
		return nil, err
	}
	receipts := []*model.Receipt{}
	for _, receipt := range receiptMap {
		user, _ := GetUserById(receipt.UserID)
		receipt := model.Receipt{
			ID:          receipt.ID,
			ImageName:   receipt.ImageName,
			ImageURL:    "http://localhost:8080/" + "image/" + receipt.ImageName,
			User:        user,
			DateCreated: receipt.DateCreated,
		}
		receipts = append(receipts, &receipt)
	}

	return receipts, nil
}

func GetReceptByID(id string) (*model.Receipt, error) {
	receiptMap, err := readReceiptMap()
	if err != nil {
		return nil, err
	}

	if receiptInternal, ok := receiptMap[id]; ok {
		user, _ := GetUserById(receiptInternal.UserID)
		receipt := model.Receipt{
			ID:          receiptInternal.ID,
			ImageName:   receiptInternal.ImageName,
			ImageURL:    "http://localhost:8080/" + "image/" + receiptInternal.ImageName,
			User:        user,
			DateCreated: receiptInternal.DateCreated,
		}
		return &receipt, nil
	}
	return nil, errors.New("No receipt with such id")
}
