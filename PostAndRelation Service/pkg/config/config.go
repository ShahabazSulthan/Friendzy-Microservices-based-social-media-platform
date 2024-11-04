package config

import "github.com/spf13/viper"

type PortManager struct {
	RunnerPort string `mapstructure:"PORTNO"`
	AuthSvcUrl string `mapstructure:"AUTH_SVC_URL"`
}

type DataBase struct {
	DBUser     string `mapstructure:"DBUSER"`
	DBHost     string `mapstructure:"DBHOST"`
	DBName     string `mapstructure:"DBNAME"`
	DBPassword string `mapstructure:"DBPASSWORD"`
	DBPort     string `mapstructure:"DBPORT"`
}

type RedisConfigs struct {
	RedisHost string `mapstructure:"REDIS_HOST"`
	RedisPort string `mapstructure:"REDIS_PORT"`
	RedisDB   int    `mapstructure:"REDIS_DB"`
}

type KafkaConfigs struct {
	KafkaPort              string `mapstructure:"KAFKA_PORT"`
	KafkaTopicNotification string `mapstructure:"KAFKA_TOPIC_2"`
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
	Kafka    KafkaConfigs
	Redis    RedisConfigs
	Smtp     Smtp
}

func LoadConfig() (*Config, error) {
	var portmanager PortManager
	var db DataBase
	var kafka KafkaConfigs
	var redis RedisConfigs
	var smtp Smtp

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

	err = viper.Unmarshal(&kafka)
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&redis)
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&smtp)
	if err != nil {
		return nil, err
	}

	// Include both PortManager and DataBase in the Config struct
	config := Config{
		PortMngr: portmanager,
		DB:       db,
		Kafka:    kafka,
		Redis:    redis,
		Smtp:     smtp,
	}

	return &config, nil
}
