package routes

import (
    "encoding/json"
    format "fmt"
    "github.com/micro/go-micro/client"
    classesAPI "github.com/plexmediamanager/micro-manager/classes/api"
    "github.com/plexmediamanager/service"
    "log"
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

func HandleAPIInformation(writer http.ResponseWriter, request *http.Request) {
    sendResponse(writer, &SuccessfulResponse {
        Message:     "Successfully fetched API information",
        Data:        classesAPI.GetAPIGeneralInformation(),
    }, http.StatusOK)
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

// Get micro client instance
func microClient() client.Client {
    if application, ok := service.FromContext(); ok {
        return application.Service().Client()
    }
    log.Panic("Well, it happened.... There was no context, no idea why.")
    return nil
}

// Build redis cache key
func buildRedisCacheKey(keyType string, tableName string, saveType interface{}) string {
    return format.Sprintf("%s::%s:%s",
        keyType,
        tableName,
        format.Sprint(saveType),
    )
}