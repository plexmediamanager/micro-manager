package routes

import (
    "encoding/json"
    format "fmt"
    classesAPI "github.com/plexmediamanager/micro-manager/classes/api"
    "net/http"
    "time"
)

type SuccessfulResponse struct {
    Success                 bool            `json:"success"`
    Data                    interface{}     `json:"data"`
    Message                 string          `json:"message"`
    RequestedOn             int64           `json:"requested_on"`
}

type FailedResponse struct {
    Success                 bool            `json:"success"`
    Data                    interface{}     `json:"data"`
    Message                 string          `json:"message"`
    RequestedOn             int64           `json:"requested_on"`
}

// Return successful response
func sendResponse(writer http.ResponseWriter, response *SuccessfulResponse, statusCode int) {
    writer.Header().Set("Content-Type", "application/json")
    writer.WriteHeader(statusCode)
    response.Success = true
    response.RequestedOn = time.Now().Unix()
    err := json.NewEncoder(writer).Encode(response)
    if err != nil {
        format.Println("Failed to encode response")
    }
}

// Return response with errors
func sendError(writer http.ResponseWriter, response *FailedResponse, statusCode int) {
    writer.Header().Set("Content-Type", "application/json")
    writer.WriteHeader(statusCode)
    response.Success = false
    response.RequestedOn = time.Now().Unix()
    err := json.NewEncoder(writer).Encode(response)
    if err != nil {
        format.Println("Failed to encode response")
    }
}

func HandleAPIInformation(writer http.ResponseWriter, request *http.Request) {
    sendResponse(writer, &SuccessfulResponse {
        Message:     "Successfully fetched API information",
        Data:        classesAPI.GetAPIGeneralInformation(),
    }, http.StatusOK)
}
