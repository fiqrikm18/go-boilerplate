package lib

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	uuid "github.com/satori/go.uuid"
	"io/ioutil"
	"os"
	"time"
)

type (
	JwtToken struct {
		verifyKey *rsa.PublicKey
		signKey   *rsa.PrivateKey
		config    *ApplicationConfig
	}

	TokenPayload struct {
		Username string
	}

	TokenData struct {
		AccessToken        string
		RefreshToken       string
		AccessTokenUUID    string
		RefreshTokenUUID   string
		AccessTokenExpire  int64
		RefreshTokenExpire int64
	}

	TokenClaims struct {
		TokenUUID    string
		TokenExpired float64
		Username     string
	}
)

func NewToken() (JwtToken, error) {
	currentDir, _ := os.Getwd()

	signBytes, err := os.ReadFile(fmt.Sprintf("%s/certs/app.rsa", currentDir))
	if err != nil {
		return JwtToken{}, err
	}

	signKey, err := jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		return JwtToken{}, err
	}

	verifyBytes, err := ioutil.ReadFile(fmt.Sprintf("%s/certs/app.rsa.pub", currentDir))
	if err != nil {
		return JwtToken{}, err
	}

	verifyKey, err := jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		return JwtToken{}, err
	}

	cfg, err := LoadConfigFile()
	if err != nil {
		return JwtToken{}, err
	}

	return JwtToken{
		signKey:   signKey,
		verifyKey: verifyKey,
		config:    cfg,
	}, nil
}

func (t *JwtToken) GenerateToken(payload TokenPayload) (*TokenData, error) {
	accessTokenExpiredDay := t.config.SecurityConf.AccessTokenExpDays
	refreshTokenExpiredDay := t.config.SecurityConf.RefreshTokenExpDays

	tokenData := &TokenData{
		AccessTokenUUID:    uuid.NewV4().String(),
		RefreshTokenUUID:   uuid.NewV4().String(),
		AccessTokenExpire:  time.Now().Add(time.Hour * 24 * time.Duration(accessTokenExpiredDay)).Unix(),
		RefreshTokenExpire: time.Now().Add(time.Hour * 24 * time.Duration(refreshTokenExpiredDay)).Unix(),
	}

	accessTokenPayload := jwt.MapClaims{}
	accessTokenPayload["username"] = payload.Username
	accessTokenPayload["uuid"] = tokenData.AccessTokenUUID
	accessTokenPayload["exp"] = tokenData.AccessTokenExpire

	accessToken := jwt.NewWithClaims(jwt.SigningMethodRS256, accessTokenPayload)
	accessTokenString, err := accessToken.SignedString(t.signKey)
	if err != nil {
		return nil, err
	}
	tokenData.AccessToken = accessTokenString

	refreshTokenPayload := jwt.MapClaims{}
	refreshTokenPayload["username"] = payload.Username
	refreshTokenPayload["uuid"] = tokenData.RefreshTokenUUID
	refreshTokenPayload["exp"] = tokenData.RefreshTokenExpire

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodRS256, refreshTokenPayload)
	refreshTokenString, err := refreshToken.SignedString(t.signKey)
	if err != nil {
		return nil, err
	}
	tokenData.RefreshToken = refreshTokenString

	return tokenData, nil
}

func (t *JwtToken) VerifyToken(tokenString string) (bool, error) {
	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.verifyKey, nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func (t *JwtToken) ExtractToken(tokenString string) (*TokenClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return t.verifyKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token")
	}

	tokenData := TokenClaims{
		TokenUUID:    claims["uuid"].(string),
		TokenExpired: claims["exp"].(float64),
		Username:     claims["username"].(string),
	}

	return &tokenData, nil
}
