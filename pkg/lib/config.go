package lib

import (
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type (
	ApplicationConfig struct {
		HttpPort int    `mapstructure:"http_port"`
		GrpcPort int    `mapstructure:"grpc_port"`
		Env      string `mapstructure:"env"`

		DBConf       DbConfig       `mapstructure:"db_conf"`
		DBConfTest   DbConfig       `mapstructure:"db_test"`
		SecurityConf SecurityConfig `mapstructure:"security_conf"`
	}

	DbConfig struct {
		Driver   string `mapstructure:"driver"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"passowrd"`
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Name     string `mapstructure:"name"`
		SSLMode  string `mapstructure:"ssl_mode"`
	}

	SecurityConfig struct {
		AccessTokenExpDays  int    `mapstructure:"access_token_exp_days"`
		RefreshTokenExpDays int    `mapstructure:"refresh_token_exp_days"`
		PasswordSalt        string `mapstructure:"password_salt"`
	}
)

func LoadConfigFile() (conf *ApplicationConfig, err error) {
	wd, err := os.Getwd()
	var path string
	if err != nil {
		return
	}

	if os.Getenv("APP_CONFIG_PATH") != "" {
		path = os.Getenv("APP_CONFIG_PATH")
	} else {
		path = filepath.Join(wd, "config")
	}

	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AutomaticEnv()
	viper.WatchConfig()

	if err = viper.ReadInConfig(); err != nil {
		return
	}

	viper.Unmarshal(&conf)
	return
}
