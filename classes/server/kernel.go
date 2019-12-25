package server

import (
    format "fmt"
    "github.com/plexmediamanager/service/log"
    "strings"
)

type KernelInformation struct {
    Version                 string      `json:"version"`
    OperatingSystem         string      `json:"os"`
    OperatingSystemVersion  string      `json:"os_version"`
    BuildDate               string      `json:"build_date"`
}

func getWindowsVersion() string {
    result, err := execute("cmd", "ver")
    if err != nil {
        log.Panic(err)
    }
    response := strings.Replace(result,"\n","",-1)
    response = strings.Replace(response,"\r\n","",-1)
    indexOfStart := strings.Index(response,"[Version")
    indexOfEnd := strings.Index(response,"]")
    var version string
    if indexOfStart == -1 || indexOfEnd == -1 {
        version = "unknown"
    } else {
        version = response[indexOfStart + 9:indexOfEnd]
    }
    return version
}

func getLinuxVersion(version *string, operatingSystem *string, operatingSystemVersion *string, buildDate *string) {
    var err error
    result, err := execute("sh",
        "-c",
        "cat /proc/version",
    )
    if err != nil {
        log.Panic(err)
    }
    partedResult := strings.Split(result, " ")
    *version = strings.TrimSpace(partedResult[2])
    *operatingSystem = strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(partedResult[7], ")", ""), "(", ""))
    indexOfStart := strings.Index(partedResult[8],"~")
    *operatingSystemVersion = strings.TrimSpace(strings.ReplaceAll(partedResult[8][indexOfStart+1:], ")", ""))
    *buildDate = format.Sprintf("%s %s, %s %s", partedResult[12], partedResult[13], strings.ReplaceAll(partedResult[16], "\n", ""), partedResult[14])
}

func getKernelInformation() *KernelInformation {
    var version string
    var operatingSystem string
    var operatingSystemVersion string
    var buildDate string

    switch operatingSystemType {
        case OS_WINDOWS:
            version = "unknown"
            operatingSystem = operatingSystemMap[operatingSystemType]
            operatingSystemVersion = getWindowsVersion()
            buildDate = "unknown"
        case OS_LINUX:
            getLinuxVersion(&version, &operatingSystem, &operatingSystemVersion, &buildDate)
        case OS_MACOS:

    }

    return &KernelInformation {
        Version:                version,
        OperatingSystem:        operatingSystem,
        OperatingSystemVersion: operatingSystemVersion,
        BuildDate:              buildDate,
    }
}