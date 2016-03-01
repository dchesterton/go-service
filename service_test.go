package service

import (
	"testing"
)

type mockService struct{}

func TestNewService(t *testing.T) {
	service := mockService{}

	svc := NewService("Service", service)

	if svc.Name != "Service" {
		t.Errorf("Expecting service name to be 'Service', got '%s'", svc.Name)
	}

	if svc.Service != service {
		t.Errorf("Expecting service to be []string{}, got '%v'", svc.Service)
	}

	if svc.isFailing {
		t.Errorf("Expected isFailing to be false by default")
	}
}
