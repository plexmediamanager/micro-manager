package main

import (
    "github.com/micro/go-micro/client"
    "github.com/plexmediamanager/micro-manager/web"
    "github.com/plexmediamanager/service"
    "github.com/plexmediamanager/service/log"
    "time"
)

func main() {
    application := service.CreateApplication()

    err := application.InitializeConfiguration()
    if err != nil {
        log.Panic(err)
    }

    err = application.InitializeMicroService()
    if err != nil {
       log.Panic(err)
    }

    err = application.Service().Client().Init(
       client.PoolSize(10),
       client.Retries(30),
       client.RequestTimeout(1 * time.Second),
    )
    if err != nil {
       log.Panic(err)
    }

    go application.StartMicroService()

    go func() {
        err := web.StartServer(application)
        if err != nil {
            log.Panic(err)
        }
    }()

    service.WaitForOSSignal(1)
}
