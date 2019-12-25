package api

import (
    "time"
)

type ApplicationInformation struct {
    LatestVersion           string                  `json:"latest_version"`
    LocalVersion            string                  `json:"local_version"`
    Version                 string                  `json:"version"`
    Updated                 bool                    `json:"updated"`
}

type ServerTime struct {
    Exact                   int64                   `json:"exact"`
    Nice                    string                  `json:"nice"`
    Timezone                string                  `json:"timezone"`
}

type GeneralInformation struct {
    Application             ApplicationInformation  `json:"application"`
    ServerTime              ServerTime              `json:"server_time"`
}

// Get general API information
// TODO: Add information about all services
func GetAPIGeneralInformation() *GeneralInformation {
    timeNow := time.Now()
    timeZone, _ := timeNow.Zone()
    timeFormat := "2006-01-02 15:04:05"
    return &GeneralInformation{
        ServerTime:     ServerTime{
            Exact:      timeNow.Unix(),
            Nice:       timeNow.Format(timeFormat),
            Timezone:   timeZone,
        },
    }
}
