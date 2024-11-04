package interface_kafka

import requestmodel "github.com/ShahabazSulthan/friendzy_post/pkg/models/requestModel"

type IKafkaProducer interface {
	KafkaNotificationProducer(message *requestmodel.KafkaNotification) error
}
