package config

import "github.com/spf13/viper"

type PortManager struct {
	RunnerPort string `mapstructure:"PORTNO"`
	AuthSvcUrl string `mapstructure:"AUTH_SVC_URL"`
}

type DataBase struct {
	DBUser     string `mapstructure:"DBUSER"`
	DBName     string `mapstructure:"DBNAME"`
	DBPassword string `mapstructure:"DBPASSWORD"`
	DBHost     string `mapstructure:"DBHOST"`
	DBPort     string `mapstructure:"DBPORT"`
}

type KafkaConfigs struct {
	KafkaPort              string `mapstructure:"KAFKA_PORT"`
	KafkaTopicNotification string `mapstructure:"KAFKA_TOPIC_1"`
}

type Config struct {
	PortMngr    PortManager
	DB          DataBase
	KafkaConfig KafkaConfigs
}

func LoadConfig() (*Config, error) {
	var portmngr PortManager
	var db DataBase
	var kafkaconfigs KafkaConfigs

	viper.AddConfigPath("./")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	err = viper.Unmarshal(&portmngr)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&db)
	if err != nil {
		return nil, err
	}
	err = viper.Unmarshal(&kafkaconfigs)
	if err != nil {
		return nil, err
	}

	config := Config{PortMngr: portmngr, DB: db, KafkaConfig: kafkaconfigs}
	return &config, nil

}
