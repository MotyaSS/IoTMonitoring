package service

import (
	"github.com/MotyaSS/IoTMonitoring/internal/scrapper/kafka/producer"
)

type Service struct {
	kafka *producer.KafkaProducer
}
