package repository

import (
	"errors"
	"github.com/fiqrikm18/go-boilerplate/pkg/lib"
	"golang.org/x/crypto/bcrypt"

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

func (repo *UserRepository) Create(userData dto.RegisterRequest) error {
	appConf, err := lib.LoadConfigFile()
	if err != nil {
		return err
	}

	userPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password+appConf.SecurityConf.PasswordSalt), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := dao.User{
		Name:     userData.Name,
		Username: userData.Username,
		Email:    userData.Email,
		Password: string(userPassword),
	}

	tx := repo.DbConn.DB.Create(&user)
	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

func (repo *UserRepository) FindByUsernameOrEmail(username string, email string) (*dao.User, error) {
	var user dao.User
	tx := repo.DbConn.DB.Where("username LIKE ? OR email LIKE ?", "%"+username+"%", "%"+email+"%").Find(&user)
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
	var user dao.User
	tx := repo.DbConn.DB.Delete(&user, id)
	if tx.Error != nil {
		return tx.Error
	}

	return nil
}
