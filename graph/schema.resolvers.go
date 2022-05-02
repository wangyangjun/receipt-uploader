package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/nfnt/resize"
	uuid "github.com/satori/go.uuid"
	"github.com/wangyangjun/receipt-uploader/graph/generated"
	"github.com/wangyangjun/receipt-uploader/graph/model"
	"github.com/wangyangjun/receipt-uploader/graph/service"
)

const dateFormat = "2006-01-02 00:00:00"

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (*model.User, error) {
	user := model.User{
		ID:          fmt.Sprintf("%v", uuid.NewV4()),
		FullName:    input.FullName,
		Email:       input.Email,
		DateCreated: time.Now().Format(dateFormat),
	}

	service.CreateUser(user)
	return &user, nil
}

func (r *mutationResolver) UploadReceipt(ctx context.Context, input model.ReceiptImage) (*model.Receipt, error) {
	user, err := service.GetUserById(input.UserID)
	if err != nil {
		fmt.Printf("Can not find user %v", err)
		return nil, errors.New("User doesn't exist")
	}
	receiptId := fmt.Sprintf("%v", uuid.NewV4())
	imageFileName := fmt.Sprintf("%v-%v", receiptId, input.File.Filename)
	receipt := model.Receipt{
		ID:          receiptId,
		ImageName:   imageFileName,
		User:        user,
		DateCreated: time.Now().Format(dateFormat),
	}

	buf := &bytes.Buffer{}
	tee := io.TeeReader(input.File.File, buf)
	// check whether it is a valid image file
	_, _, err = image.Decode(tee)
	if err != nil {
		return nil, errors.New("Unsupported file type, only png, jpg and gif are supported")
	}

	stream, err := ioutil.ReadAll(buf)
	if err != nil {
		fmt.Printf("error from file %v", err)
	}
	fileErr := ioutil.WriteFile("images/"+imageFileName, stream, 0644)
	if fileErr != nil {
		fmt.Printf("file err %v", fileErr)
	}

	service.CreateRecept(receipt)

	return &receipt, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	return service.GetAllUsers()
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return service.GetUserById(id)
}

func (r *queryResolver) Receipts(ctx context.Context) ([]*model.Receipt, error) {
	return service.GetAllRecepts()
}

func (r *queryResolver) ReceiptImage(ctx context.Context, id string, resolution *int) (*model.Receipt, error) {
	receipt, err := service.GetReceptByID(id)

	if resolution != nil && *resolution != 100 {
		imageFileNameWithScale := strconv.Itoa(*resolution) + "-" + receipt.ImageName
		_, err := os.Stat("images/" + imageFileNameWithScale)

		if errors.Is(err, os.ErrNotExist) {
			_, err := os.Stat("images/" + receipt.ImageName)
			if errors.Is(err, os.ErrNotExist) {
				return nil, errors.New("Image cannot be found for the receipt")
			}
			file, err := os.Open("images/" + receipt.ImageName)
			defer file.Close()
			log.Println("images/" + receipt.ImageName)

			if err != nil {
				return nil, errors.New("Image cannot be opened")
			}

			img, _, err := image.Decode(file)
			if err != nil {
				return nil, err
			}

			file.Seek(0, 0)
			imgConfig, _, err := image.DecodeConfig(file)
			if err != nil {
				return nil, errors.New("Decode image config failed")
			}
			newImage := resize.Resize(uint(float32(imgConfig.Width*(*resolution))*0.01), 0, img, resize.Lanczos3)
			scaleImageFile, err := os.Create("images/" + imageFileNameWithScale)

			if err != nil {
				log.Fatal(err)
			}
			defer scaleImageFile.Close()

			// write new image to file
			jpeg.Encode(scaleImageFile, newImage, nil)

		}

	}

	return receipt, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
