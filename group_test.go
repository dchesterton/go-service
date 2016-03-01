package service

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type ServiceA struct{}

func TestGetFailingServices(t *testing.T) {
	serviceA := NewService("Service A", ServiceA{})
	serviceA.isFailing = true

	serviceB := NewService("Service B", ServiceA{})
	serviceC := NewService("Service C", ServiceA{})

	group := NewServiceGroup(serviceA, serviceB, serviceC)

	failing := group.getFailing()

	if len(failing) != 1 {
		t.Errorf("Expected 1 failing service, got %d", len(failing))
	}

	nonFailing := group.getNonFailing()

	if len(nonFailing) != 2 {
		t.Errorf("Expected 2 non failing services, got %d", len(nonFailing))
	}
}

func TestResetFailingServices(t *testing.T) {
	failureTime := time.Now()
	failureTime = failureTime.Add(-(time.Duration(5) * time.Minute))

	serviceA := NewService("Service A", ServiceA{})
	serviceA.isFailing = true
	serviceA.lastFailure = failureTime

	serviceB := NewService("Service B", ServiceA{})

	group := NewServiceGroup(serviceA, serviceB)
	group.resetFailing()

	failing := group.getFailing()

	if len(failing) != 0 {
		t.Errorf("Expected 0 failing services, got %d", len(failing))
	}
}

func TestDoesNotResetRecentFailingServices(t *testing.T) {
	failureTime := time.Now()
	failureTime = failureTime.Add(-(time.Duration(2) * time.Minute))

	serviceA := NewService("Service A", ServiceA{})
	serviceA.isFailing = true
	serviceA.lastFailure = failureTime

	serviceB := NewService("Service B", ServiceA{})

	group := NewServiceGroup(serviceA, serviceB)
	group.resetFailing()

	failing := group.getFailing()

	if len(failing) != 1 {
		t.Errorf("Expected 1 failing service, got %d", len(failing))
	}
}

func TestTryCallsNextServiceOnError(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	serviceA := NewService("Service A", ServiceA{})
	serviceB := NewService("Service B", ServiceA{})
	serviceC := NewService("Service C", ServiceA{})

	group := NewServiceGroup(serviceA, serviceB, serviceC)
	group.LoadBalance = true

	i := 0

	serviceErr := group.Try(func(service *Service) error {
		var err error

		// error for first two services
		if i <= 1 {
			err = fmt.Errorf("An error")
		}

		i = i + 1

		return err
	})

	if i != 3 {
		t.Errorf("Expected Try to be called 3 times, called %d times", i)
	}

	if serviceErr != nil {
		t.Errorf("Not expecting error, got %v", serviceErr)
	}
}

func TestTryReturnsError(t *testing.T) {
	serviceA := NewService("Service A", ServiceA{})
	serviceB := NewService("Service B", ServiceA{})
	serviceC := NewService("Service C", ServiceA{})

	group := NewServiceGroup(serviceA, serviceB, serviceC)

	i := 0

	serviceErr := group.Try(func(service *Service) error {
		i = i + 1
		return fmt.Errorf("Error %d", i)
	})

	if serviceErr == nil {
		t.Error("Expecting an error, got nil")
	}

	if serviceErr.Error() != "3 services failed, could not complete action (Error 1, Error 2, Error 3)" {
		t.Errorf("Unexpected error message, got %v", serviceErr.Error())
	}

	if !serviceA.isFailing {
		t.Error("Service A should be marked as failing")
	}

	if !serviceB.isFailing {
		t.Error("Service B should be marked as failing")
	}

	if !serviceC.isFailing {
		t.Error("Service C should be marked as failing")
	}
}

func TestFailingServiceIsCalledLast(t *testing.T) {
	serviceA := NewService("Service A", ServiceA{})
	serviceA.isFailing = true
	serviceA.lastFailure = time.Now()

	serviceB := NewService("Service B", ServiceA{})

	group := NewServiceGroup(serviceA, serviceB)

	var first *Service

	group.Try(func(service *Service) error {
		first = service
		return nil
	})

	if first != serviceB {
		t.Error("Expected service B to be called first")
	}
}

func TestPreviouslyFailingServiceIsUnmarked(t *testing.T) {
	serviceA := NewService("Service A", ServiceA{})
	serviceA.isFailing = true
	serviceA.lastFailure = time.Now()

	group := NewServiceGroup(serviceA)

	group.Try(func(service *Service) error {
		return nil
	})

	if serviceA.isFailing {
		t.Error("Expected service to no longer be failing")
	}
}
