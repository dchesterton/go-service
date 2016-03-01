package service

import (
	"time"
)

type Service struct {
	Name        string
	Service     interface{}
	isFailing   bool
	lastFailure time.Time
}

func NewService(name string, service interface{}) *Service {
	return &Service{
		Name:      name,
		Service:   service,
		isFailing: false,
	}
}
