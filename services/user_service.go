package services

import "example.com/tes-rmq/models"

type UserService interface {
	Register(*models.User) error
}