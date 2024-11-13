package config

import "github.com/spf13/viper"

type PortManager struct {
	PortNo         string `mapstructure:"PORTNO"`
	PostNrelSvcUrl string `mapstructure:"POSTNREL_SVC_URL"`
}

type Razopay struct {
	RazopayKey    string `mapstructure:"RAZOPAYKEY"`
	RazopaySecret string `mapstructure:"PAZOPAYSECRET"`
}

type DataBase struct {
	DBUser     string `mapstructure:"DBUSER"`
	DBHost     string `mapstructure:"DBHOST"`
	DBName     string `mapstructure:"DBNAME"`
	DBPassword string `mapstructure:"DBPASSWORD"`
	DBPort     string `mapstructure:"DBPORT"`
}

type Token struct {
	AdminSecurityKey    string `mapstructure:"ADMIN_TOKENKEY"`
	UserSecurityKey     string `mapstructure:"USER_TOKENKEY"`
	TempVerificationKey string `mapstructure:"TEMPERVERY_TOKENKEY"`
}

type Smtp struct {
	SmtpSender   string `mapstructure:"FROM_EMAIL"`
	SmtpPassword string `mapstructure:"SMTP_PASSWORD"`
	SmtpHost     string `mapstructure:"SMTP_SERVER"`
	SmtpPort     string `mapstructure:"SMTP_PORT"`
}

type Config struct {
	PortMngr PortManager
	DB       DataBase
	Token    Token
	Smtp     Smtp
	Razopay  Razopay
}

func LoadConfig() (*Config, error) {
	var portmanager PortManager
	var token Token
	var smtp Smtp
	var db DataBase
	var razoPay Razopay

	viper.AddConfigPath("./")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&portmanager)
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&db)
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&token)
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&smtp)
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&razoPay)
	if err != nil {
		return nil, err
	}

	config := Config{PortMngr: portmanager, Token: token, DB: db, Smtp: smtp, Razopay: razoPay}
	return &config, nil
}
