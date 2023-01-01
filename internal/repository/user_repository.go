package repository

import (
	"errors"

	"github.com/fiqrikm18/go-boilerplate/internal/config"
	"github.com/fiqrikm18/go-boilerplate/internal/model/dao"
	"github.com/fiqrikm18/go-boilerplate/internal/model/dto"
)

type UserRepository struct {
	DbConn *config.DbConnection
}

func NewUserRepository() (*UserRepository, error) {
	conn, err := config.NewDbConnection()
	if err != nil {
		return nil, err
	}

	if conn == nil {
		return nil, errors.New("connection to database failed")
	}

	return &UserRepository{DbConn: conn}, nil
}

func (repo *UserRepository) Create(userData *dto.UserRequest) error {
	return nil
}

func (repo *UserRepository) FindByUsernameOrEmail(keyword string) (*dao.User, error) {
	var user dao.User
	tx := repo.DbConn.DB.Where("username LIKE ? OR email LIKE ?", "%"+keyword+"%", "%"+keyword+"%").Find(&user)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}

func (repo *UserRepository) FindByUserId(id int) (*dao.User, error) {
	var user dao.User
	tx := repo.DbConn.DB.Find(&user, id)
	if tx.Error != nil {
		return nil, tx.Error
	}

	return &user, nil
}

func (repo *UserRepository) Delete(id int) error {
	return nil
}
