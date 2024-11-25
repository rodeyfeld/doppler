package user

import (
	"github.com/rodeyfeld/doppler/internal/requests"

	"gorm.io/gorm"
)

type ServiceWrapper interface {
	Register(request *requests.RegisterRequest) error
}

type Service struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) *Service {
	return &Service{DB: db}
}
