package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/satori/go.uuid"
	"github.com/wangyangjun/receipt-uploader/graph/generated"
	"github.com/wangyangjun/receipt-uploader/graph/model"
	"github.com/wangyangjun/receipt-uploader/graph/service"
	"github.com/wangyangjun/receipt-uploader/graph/service/auth"
)

func (r *mutationResolver) Signup(ctx context.Context, username string, password string) (*model.User, error) {
	user, err := service.CreateUser(username, password)
	return user, err
}

func (r *mutationResolver) Login(ctx context.Context, username string, password string) (*model.AuthPayload, error) {
	return service.Login(username, password)
}

func (r *mutationResolver) UploadReceipt(ctx context.Context, description string, file graphql.Upload) (*model.Receipt, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("Access denied")
	}

	receiptId := fmt.Sprintf("%v", uuid.NewV4())
	imageFileName := fmt.Sprintf("%v-%v", receiptId, file.Filename)
	receipt := model.Receipt{
		ID:          receiptId,
		ImageName:   imageFileName,
		Description: description,
		DateCreated: time.Now().Format(dateFormat),
	}

	err := service.SaveReceiptImg(imageFileName, file, user.ID)
	if err != nil {
		return nil, err
	}
	err = service.CreateRecept(receipt, user.ID)
	if err != nil {
		return nil, err
	}
	return &receipt, nil
}

func (r *queryResolver) Users(ctx context.Context) ([]*model.User, error) {
	return service.GetAllUsers()
}

func (r *queryResolver) Receipts(ctx context.Context) ([]*model.Receipt, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("Access denied")
	}
	return service.GetAllReceipts(user.ID)
}

func (r *queryResolver) Receipt(ctx context.Context, id string, scaleRatio *int) (*model.Receipt, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("Access denied")
	}
	receipt, err := service.GetReceptByID(id, user.ID)
	if err != nil {
		return nil, err
	}
	if scaleRatio != nil && *scaleRatio != 100 {
		if *scaleRatio <= 0 || *scaleRatio > 100 {
			return nil, errors.New("Not a valid scale ratio")
		}
		imageFileNameWithScale := strconv.Itoa(*scaleRatio) + "-" + receipt.ImageName
		_, err := os.Stat("image/" + imageFileNameWithScale)

		if errors.Is(err, os.ErrNotExist) {
			err = service.ScaleReceiptImage(receipt, *scaleRatio, user.ID)
			if err != nil {
				return nil, err
			}
		}
		receipt.ImageURL = service.ImageUrl(user.ID, imageFileNameWithScale)
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
