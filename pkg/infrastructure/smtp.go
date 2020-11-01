package infrastructure

import (
	"github.com/spf13/viper"
)

type SmtpConf struct {
	Host     string
	Port     string
	Login    string
	Password string
}

func InitSmtp(path string) (*SmtpConf, error) {
	viper.SetConfigFile(path)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	return &SmtpConf{
		Host:     viper.GetString("host"),
		Port:     viper.GetString("port"),
		Login:    viper.GetString("login"),
		Password: viper.GetString("password"),
	}, nil
}
