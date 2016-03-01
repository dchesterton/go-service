package service

import (
	"fmt"
	"strings"
	"time"
)

type Logger interface {
	Log(err error)
}

type ServiceGroup struct {
	Services    serviceList
	Logger      Logger
	RetryDelay  time.Duration
	LoadBalance bool
}

func NewServiceGroup(services ...*Service) *ServiceGroup {
	return &ServiceGroup{
		Services:    serviceList(services),
		RetryDelay:  time.Duration(5) * time.Minute, // wait 5 minutes by default
		LoadBalance: false,
	}
}

/**
 * Get all services which have not failed recently.
 */
func (group *ServiceGroup) getNonFailing() serviceList {
	services := serviceList{}

	for _, service := range group.Services {
		if !service.isFailing {
			services = append(services, service)
		}
	}

	// if we're load balancing, shuffle the array
	if group.LoadBalance {
		services.Shuffle()
	}

	return services
}

/**
 * Get all services which have failed recently.
 */
func (group *ServiceGroup) getFailing() serviceList {
	services := serviceList{}

	for _, service := range group.Services {
		if service.isFailing {
			services = append(services, service)
		}
	}

	return services
}

/**
 * Reset any services which have previously failed but delay
 * has elapsed.
 */
func (group *ServiceGroup) resetFailing() {
	for _, service := range group.Services {
		if service.isFailing {
			if time.Since(service.lastFailure).Seconds() > group.RetryDelay.Seconds() {
				service.isFailing = false
			}
		}
	}
}

/**
 * Try and execute service.
 */
func (group *ServiceGroup) Try(fn func(service *Service) error) error {
	// reset any failing services
	group.resetFailing()

	nonFailingServices := group.getNonFailing()
	failingServices := group.getFailing()

	allServices := serviceList{}

	for _, service := range nonFailingServices {
		allServices = append(allServices, service)
	}

	for _, service := range failingServices {
		allServices = append(allServices, service)
	}

	errorStrings := []string{}

	for _, service := range allServices {
		err := fn(service)

		if err == nil {
			service.isFailing = false
			return nil
		} else {
			// mark the service as failing which will ensure it's tried last on the next run to stop
			// potential overloading issues
			service.isFailing = true
			service.lastFailure = time.Now()

			errorStrings = append(errorStrings, err.Error())

			//group.Logger.Log(err)
		}
	}

	errors := strings.Join(errorStrings, ", ")

	return fmt.Errorf("%d services failed, could not complete action (%s)", len(group.Services), errors)
}
