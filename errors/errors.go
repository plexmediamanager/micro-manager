
package errors

import "github.com/plexmediamanager/service/errors"

const (
    ServiceID       errors.Service      =   6
)

var (
    // Library errors
    UnableToCreateHTTPServer            =   errors.Error {
        Code:               errors.Code {
            Service:        ServiceID,
            ErrorType:      errors.TypeLibrary,
            ErrorNumber:    1,
        },
        Message:            "Unable to start server at: %s",
    }
    MarshalError                        =   errors.Error {
        Code:               errors.Code {
            Service:        ServiceID,
            ErrorType:      errors.TypeLibrary,
            ErrorNumber:    2,
        },
        Message:            "Unable to marshal structure",
    }
    UnmarshalError                      =   errors.Error {
        Code:               errors.Code {
            Service:        ServiceID,
            ErrorType:      errors.TypeLibrary,
            ErrorNumber:    3,
        },
        Message:            "Unable to unmarshal structure",
    }

    // Service errors
    RedisGetError                       =   errors.Error {
        Code:               errors.Code {
            Service:        ServiceID,
            ErrorType:      errors.TypeService,
            ErrorNumber:    1,
        },
        Message:            "Unable to get value from Redis: %s",
    }
    RedisSetError                       =   errors.Error {
        Code:               errors.Code {
            Service:        ServiceID,
            ErrorType:      errors.TypeService,
            ErrorNumber:    2,
        },
        Message:            "Unable set value to Redis: %s",
    }

)