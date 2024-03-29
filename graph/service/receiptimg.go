package service

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/nfnt/resize"
	"github.com/wangyangjun/receipt-uploader/graph/model"
)

// validate upload imge file and save it
func SaveReceiptImg(imageFileName string, file graphql.Upload, userId string) error {

	if _, err := os.Stat(path.Join("images", userId)); os.IsNotExist(err) {
		err := os.MkdirAll(path.Join("images", userId), 0700)
		if err != nil {
			log.Fatal("image directory for user initialized failed!")
		}
	}

	buff, err := ioutil.ReadAll(file.File)
	if err != nil {
		panic(err)
	}

	// check where the upload file is a valid image
	reader1 := bytes.NewReader(buff)
	_, _, err = image.DecodeConfig(reader1)
	if err != nil {
		return errors.New("Unsupported file type, only png, jpg and gif are supported")
	}

	// save image on local file system
	reader2 := bytes.NewReader(buff)
	stream, err := ioutil.ReadAll(reader2)
	if err != nil {
		log.Printf("error from file %v", err)
		return err
	}
	err = ioutil.WriteFile(path.Join("images", userId, imageFileName), stream, 0644)
	if err != nil {
		log.Printf("file err %v", err)
		return err
	}
	return nil
}

// scale image by percentage and save it
func ScaleReceiptImage(receipt *model.Receipt, scaleRatio int, userId string) error {
	imageFileNameWithScale := strconv.Itoa(scaleRatio) + "-" + receipt.ImageName
	filePath := path.Join("images", userId, receipt.ImageName)

	_, err := os.Stat(filePath)
	if errors.Is(err, os.ErrNotExist) {
		return errors.New("Image cannot be found for the receipt")
	}
	file, err := os.Open(filePath)
	defer file.Close()

	if err != nil {
		return errors.New("Image cannot be opened")
	}

	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	file.Seek(0, 0)
	imgConfig, _, err := image.DecodeConfig(file)
	if err != nil {
		return errors.New("Decode image config failed")
	}
	newImage := resize.Resize(uint(float32(imgConfig.Width*scaleRatio)*0.01), 0, img, resize.Lanczos3)
	scaleImageFile, err := os.Create(path.Join("images", userId, imageFileNameWithScale))

	if err != nil {
		log.Fatal(err)
	}
	defer scaleImageFile.Close()

	// write new image to file
	jpeg.Encode(scaleImageFile, newImage, nil)
	return nil
}
