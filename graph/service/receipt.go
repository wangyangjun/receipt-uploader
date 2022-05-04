package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/wangyangjun/receipt-uploader/graph/model"
)

func CreateRecept(receipt model.Receipt, userId string) error {
	file, err := ioutil.ReadFile(receiptFilePath)
	if err != nil {
		log.Println(err)
		return errors.New("Can not open receipt file")
	}

	receiptInternal := model.ReceiptInternal{
		ID:          receipt.ID,
		ImageName:   receipt.ImageName,
		UserID:      userId,
		Description: receipt.Description,
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

func ImageUrl(userId string, imageName string) string {
	return "http://localhost:" + os.Getenv("PORT") + "/image/" + userId + "/" + imageName
}

func GetAllReceipts(userId string) ([]*model.Receipt, error) {
	receiptMap, err := readReceiptMap()
	if err != nil {
		return nil, err
	}
	receipts := []*model.Receipt{}
	for _, receipt := range receiptMap {
		if receipt.UserID == userId {
			receipt := model.Receipt{
				ID:          receipt.ID,
				ImageName:   receipt.ImageName,
				ImageURL:    ImageUrl(userId, receipt.ImageName),
				Description: receipt.Description,
				DateCreated: receipt.DateCreated,
			}
			receipts = append(receipts, &receipt)
		}
	}

	return receipts, nil
}

func GetReceptByID(id string, userId string) (*model.Receipt, error) {
	receiptMap, err := readReceiptMap()
	if err != nil {
		return nil, err
	}

	receiptInternal, ok := receiptMap[id]
	if ok && receiptInternal.UserID == userId {
		receipt := model.Receipt{
			ID:          receiptInternal.ID,
			ImageName:   receiptInternal.ImageName,
			ImageURL:    ImageUrl(userId, receiptInternal.ImageName),
			Description: receiptInternal.Description,
			DateCreated: receiptInternal.DateCreated,
		}
		return &receipt, nil
	}
	return nil, fmt.Errorf("No receipt with such id")
}
