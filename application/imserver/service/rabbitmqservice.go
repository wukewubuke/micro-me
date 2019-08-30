package service

import (
	"fmt"
	"github.com/micro/go-micro/broker"
	"log"
)

type (
	RabbitMqService struct {
		topic          string
		rabbitMqBroker broker.Broker
	}
)

func NewRabbitMqService(topic string, rabbitMqBroker broker.Broker) (*RabbitMqService, error) {
	// 初始化
	if err := rabbitMqBroker.Init(); err != nil {
		return nil, err
	}
	if err := rabbitMqBroker.Connect(); err != nil {
		return nil, err
	}

	return &RabbitMqService{topic: topic, rabbitMqBroker: rabbitMqBroker}, nil
}

func (s *RabbitMqService) Publisher(msg *broker.Message) {
	if err := s.rabbitMqBroker.Publish(s.topic, msg); err != nil {
		log.Printf("[publisher err] : %+v", err)
		return
	}
	fmt.Printf("[Publisher success: %s:%+v]\n", s.topic, msg)
}

func (s *RabbitMqService) Subscribe(f func(msg []byte) error) {
	if _, err := s.rabbitMqBroker.Subscribe(s.topic, func(event broker.Event) error {
		if err := f(event.Message().Body); err != nil {
			log.Printf("Subscribe f msg err :%+v", err)
			return err
		}
		return nil
	}); err != nil {
		log.Printf("[Subscribe %s err] : %+v", s.topic, err)
	}
}
