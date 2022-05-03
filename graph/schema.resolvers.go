package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	_ "image/gif"
	"os"
	"strconv"
	"time"

	"github.com/satori/go.uuid"
	"github.com/wangyangjun/receipt-uploader/graph/generated"
	"github.com/wangyangjun/receipt-uploader/graph/model"
	"github.com/wangyangjun/receipt-uploader/graph/service"
)

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

	err = service.SaveReceiptImg(imageFileName, input)
	if err != nil {
		return nil, err
	}
	err = service.CreateRecept(receipt)
	if err != nil {
		return nil, err
	}
	return &receipt, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	return service.GetAllUsers()
}

func (r *queryResolver) User(ctx context.Context, id string) (*model.User, error) {
	return service.GetUserById(id)
}

func (r *queryResolver) Receipts(ctx context.Context) ([]*model.Receipt, error) {
	return service.GetAllReceipts()
}

func (r *queryResolver) Receipt(ctx context.Context, id string, scaleRatio *int) (*model.Receipt, error) {
	receipt, err := service.GetReceptByID(id)

	if scaleRatio != nil && *scaleRatio != 100 {
		if *scaleRatio <= 0 || *scaleRatio > 100 {
			return nil, errors.New("Not a valid scale ratio")
		}
		imageFileNameWithScale := strconv.Itoa(*scaleRatio) + "-" + receipt.ImageName
		_, err := os.Stat("image/" + imageFileNameWithScale)

		if errors.Is(err, os.ErrNotExist) {
			err = service.ScaleReceiptImage(receipt, *scaleRatio)
			if err != nil {
				return nil, err
			}
		}
		receipt.ImageURL = "http://localhost:8080/" + "image/" + imageFileNameWithScale
	}

	return receipt, err
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

// !!! WARNING !!!
// The code below was going to be deleted when updating resolvers. It has been copied here so you have
// one last chance to move it out of harms way if you want. There are two reasons this happens:
//  - When renaming or deleting a resolver the old code will be put in here. You can safely delete
//    it when you're done.
//  - You have helper methods in this file. Move them out to keep these resolver files clean.
const dateFormat = "2006-01-02 00:00:00"
