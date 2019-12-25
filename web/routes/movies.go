package routes

import (
    databaseService "github.com/plexmediamanager/micro-database/service"
    "github.com/plexmediamanager/service/errors"
    "net/http"
)

func HandleMoviesList(writer http.ResponseWriter, request *http.Request) {
    response, err := databaseService.MovieServiceFindDownloaded(microClient())
    if err != nil {
        sendError(writer, &FailedResponse {
            Message:     "Failed to fetch list of movies",
            Data:        errors.ParseError(err),
        }, http.StatusOK)
    } else {
        sendResponse(writer, &SuccessfulResponse {
            Message:     "Successfully fetched list of movies",
            Data:        response,
        }, http.StatusOK)
    }
}
