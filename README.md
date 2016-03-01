# Go Service
A small Go (golang) library for handling failing services and load balancing across services.

## Using backups

Firstly, create a service group for any related services. Then provide a function which accepts a service as an argument.

The library will try to use the services in the order that they're given. It keeps track of any services which have
recently failed, and not call them for a while to stop the service being overloaded.

```go
primaryService := ...
backupService := ...
backupToBackupService := ...

group := service.NewServiceGroup(
    service.NewService("primary", primaryService),
    service.NewService("backup", backupService),
    service.NewService("backup_to_backup", backupToBackupService),
)

err := group.Try(func(service *Service) error {
    // do something with service

    // you should return an error if the service is not working
    // return fmt.Errorf("")

    return nil
})

if err != nil {
    fmt.Println("Could not connect to any service!")
}

```

## Load balancing

You can tell the library to load balance across different services. This is helpful if you have two external providers and
wish to split the traffic equally.

```go
firstService := ...
secondService := ...

group := service.NewServiceGroup(
    service.NewService("first", firstService),
    service.NewService("second", secondService),
)

// load balance equally across the two services
group.LoadBalance = true

err := group.Try(func(service *Service) error {
    // do something with service

    // you should return an error if the service is not working
    // return fmt.Errorf("")

    return nil
})

if err != nil {
    fmt.Println("Could not connect to first or second service!")
}

```
