package service

import (
	"errors"
	"log"
	"os"
)

const userFilePath = "data/users.txt"
const receiptFilePath = "data/receipts.txt"

func init() {
	log.Println("initialize service package")
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		err := os.Mkdir("data", 0700)
		if err != nil {
			log.Fatal("data directory initialized failed!")
		}
	}
	if _, err := os.Stat("images"); os.IsNotExist(err) {
		err := os.Mkdir("images", 0700)
		if err != nil {
			log.Fatal("images directory initialized failed!")
		}
	}

	_, err := os.Stat(userFilePath)

	if errors.Is(err, os.ErrNotExist) {
		log.Println("users file does not exist")
		userFile, err := os.Create(userFilePath)
		defer userFile.Close()

		if err != nil {
			log.Fatal("users file initialized failed!")
		}
	}
	_, err = os.Stat(receiptFilePath)
	if errors.Is(err, os.ErrNotExist) {
		log.Println("receipt file does not exist")
		userFile, err := os.Create(receiptFilePath)
		defer userFile.Close()

		if err != nil {
			log.Fatal("receipt file initialized failed!")
		}
	}
}
