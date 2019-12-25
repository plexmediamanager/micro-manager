package routes

import (
"net/http"
)

func HandleSeriesList(writer http.ResponseWriter, request *http.Request) {
    sendResponse(writer, &SuccessfulResponse {
        Message:     "Successfully fetched list of series",
        Data:        nil,
    }, http.StatusOK)
}

