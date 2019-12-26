package main

import (
    "encoding/json"
    "github.com/micro/go-micro/client"
    "github.com/plexmediamanager/micro-manager/errors"
    "github.com/plexmediamanager/micro-manager/web"
    redisService "github.com/plexmediamanager/micro-redis/service"
    tmdbService "github.com/plexmediamanager/micro-tmdb/service"
    "github.com/plexmediamanager/micro-tmdb/tmdb"
    "github.com/plexmediamanager/service"
    "github.com/plexmediamanager/service/ctx"
    "github.com/plexmediamanager/service/helpers"
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

    err = loadTMDBAPIConfiguration(application)
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

func loadTMDBAPIConfiguration(application *service.Application) error {
    microClient := application.Service().Client()
    redisKey := "tmdb::configuration:api"
    has, err := redisService.RedisHas(microClient, redisKey)
    if err != nil {
        return errors.RedisGetError.ToErrorWithArguments(err, redisKey)
    }

    if !has {
        configuration, err := tmdbService.TMDBServiceGetAPIConfiguration(microClient)
        if err != nil {
            return err
        }
        serialized, err := json.Marshal(configuration)
        if err != nil {
            return errors.MarshalError.ToError(err)
        }
        _, err = redisService.RedisSetValue(
            microClient,
            redisKey,
            string(serialized),
            helpers.HoursToDuration(24),
        )
        if err != nil {
            return errors.RedisSetError.ToErrorWithArguments(err, redisKey)
        }
        ctx.WithValue("tmdbConfiguration", configuration)
    } else {
        bytes, err := redisService.RedisGetValue(microClient, redisKey)
        if err != nil {
            return errors.RedisGetError.ToErrorWithArguments(err, redisKey)
        }
        var configuration *tmdb.APIConfiguration
        err = json.Unmarshal(bytes.Response, &configuration)
        if err != nil {
            return errors.UnmarshalError.ToError(err)
        }
        ctx.WithValue("tmdbConfiguration", configuration)
    }

    return nil
}