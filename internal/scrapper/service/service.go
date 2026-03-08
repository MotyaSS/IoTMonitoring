package service

import "github.com/segmentio/kafka-go"

type Service struct {
	kafka kafka.ProducerSession
}
