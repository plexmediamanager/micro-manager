package web

import (
    format "fmt"
    "github.com/gorilla/mux"
    "github.com/plexmediamanager/micro-manager/errors"
    "github.com/plexmediamanager/micro-manager/web/routes"
    "github.com/plexmediamanager/service"
    "github.com/plexmediamanager/service/helpers"
    "net/http"
)

var (
    application *service.Application
)

// Get address on which server will be listening for connections
func GetServerAddress() string {
    return format.Sprintf(
        "%s:%d",
        helpers.GetEnvironmentVariableAsString("SERVER_HOST", "127.0.0.1"),
        helpers.GetEnvironmentVariableAsInteger("SERVER_PORT", 8080),
    )
}

// Start Web Server
func StartServer(app *service.Application) error {
    application = app
    router := mux.NewRouter()
    //router.Handle("/*", http.NotFoundHandler())
    registerRoutes(router)
    http.Handle("/", router)
    err := http.ListenAndServe(GetServerAddress(), nil)
    if err != nil {
        return errors.UnableToCreateHTTPServer.ToErrorWithArguments(err, GetServerAddress())
    }
    return nil
}

// Register all routes
func registerRoutes(router *mux.Router) {
    router.HandleFunc("/", routes.HandleAPIInformation).Methods(http.MethodGet)
    registerDashboardRoutes(router)
}

func registerDashboardRoutes(router *mux.Router) {
    dashboard := router.PathPrefix("/dashboard").Subrouter()
    dashboard.HandleFunc("/server-information", routes.HandleDashboardServerInformation).Methods(http.MethodGet)
}