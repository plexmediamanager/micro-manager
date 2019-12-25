
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
)