package seeder

import (
	"context"
	"github.com/bxcodec/faker/v4"
	"github.com/google/uuid"
	"go_gin/internal/domain/model"
	"gorm.io/gorm"
)

type UserSeed struct {
	Repository model.UsersRepository
}

func NewUserSeed(repository model.UsersRepository) model.CustomSeeds {
	return &UserSeed{Repository: repository}
}

func (u *UserSeed) Run(ctx context.Context, DB *gorm.DB, fill int) error {
	var users model.UsersRequests
	for i := 0; i < fill; i++ {
		user := model.UserRequest{
			ID:       uuid.New(),
			Username: faker.Username(),
			Email:    faker.Email(),
			Password: faker.PASSWORD,
			Roles:    model.Basic,
		}
		users = append(users, user)
	}
	err := u.Repository.CreateUsers(ctx, DB, users.ToUsers())
	return err
}
