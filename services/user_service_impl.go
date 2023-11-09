package services

import (
	"context"

	"example.com/tes-rmq/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserServiceImpl struct {
	usercollection *mongo.Collection
	ctx context.Context
}

func NewUserService(usercollection *mongo.Collection, ctx context.Context) UserService {
	return &UserServiceImpl{
		usercollection: usercollection,
		ctx: ctx,
	}
}

func (u *UserServiceImpl) Register(user *models.User) error {


	
	_, err := u.usercollection.InsertOne(u.ctx, user)

	return err
}