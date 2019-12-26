package common

import (
    "github.com/micro/go-micro/client"
    "github.com/plexmediamanager/service"
    "github.com/plexmediamanager/service/log"
)

// Get micro client instance
func microClient() client.Client {
    if application, ok := service.FromContext(); ok {
        return application.Service().Client()
    }
    log.Panic("Well, it happened.... There was no context, no idea why.")
    return nil
}
