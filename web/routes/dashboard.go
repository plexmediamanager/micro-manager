package routes

import (
    classesServer "github.com/plexmediamanager/micro-manager/classes/server"
    "net/http"
)

func HandleDashboardServerInformation(writer http.ResponseWriter, request *http.Request) {
    sendResponse(writer, &SuccessfulResponse {
        Message:     "Successfully fetched API information",
        Data:        classesServer.GetServerInformation(),
    }, http.StatusOK)
}
