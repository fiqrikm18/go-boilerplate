package repository

import (
	"errors"
	"github.com/fiqrikm18/go-boilerplate/internal/config"
	"github.com/fiqrikm18/go-boilerplate/internal/model/dao"
)

type OAuthAccessTokenRepository struct {
	DbConn *config.DbConnection
}

func NewOAuthAccessTokenRepository() (*OAuthAccessTokenRepository, error) {
	conn, err := config.NewDbConnection()
	if err != nil {
		return nil, err
	}

	if conn == nil {
		return nil, errors.New("connection to database failed")
	}

	return &OAuthAccessTokenRepository{DbConn: conn}, nil
}

func (repo *OAuthAccessTokenRepository) Create(data dao.OauthAccessToken) error {
	tx := repo.DbConn.DB.Create(&data)
	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

func (repo *OAuthAccessTokenRepository) RevokeByUserID(userId int64) error {
	tx := repo.DbConn.DB.Model(&dao.OauthAccessToken{}).Where("user_id=?", userId).Update("revoked", true)
	if err := tx.Error; err != nil {
		return err
	}

	return nil
}
