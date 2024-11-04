package kafka

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/IBM/sarama"

	"github.com/ShahabazSulthan/friendzy_post/pkg/config"
	requestmodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/requestModel"
	interface_kafka "github.com/ShahabazSulthan/friendzy_post/pkg/utils/Kakfa/interface"
)

type KafkaProducer struct {
	Config config.KafkaConfigs
}

func NewKafkaProducer(config config.KafkaConfigs) interface_kafka.IKafkaProducer {
	return &KafkaProducer{
		Config: config,
	}
}

func (k *KafkaProducer) KafkaNotificationProducer(message *requestmodel.KafkaNotification) error {

	fmt.Println("---------------to kafkaProducer:", *message)
	configs := sarama.NewConfig()
	configs.Producer.Return.Successes = true
	configs.Producer.Retry.Max = 5

	producer, err := sarama.NewSyncProducer([]string{k.Config.KafkaPort}, configs)
	if err != nil {
		log.Println("---------kafka producer err--------", err)
		return err
	}

	msgJson, _ := marshalStructJson(message)

	msg := &sarama.ProducerMessage{Topic: k.Config.KafkaTopicNotification,
		Key:   sarama.StringEncoder(message.UserID),
		Value: sarama.StringEncoder(*msgJson)}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Printf("\nerr sending message to kafkaproducer on partition: %s,error: %v", k.Config.KafkaTopicNotification, err)
	}
	log.Printf("[producer] partition id: %d; offset:%d, value: %v\n", partition, offset, msg)
	return nil
}

func marshalStructJson(msgModel interface{}) (*[]byte, error) {
	data, err := json.Marshal(msgModel)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
