package server

import (
    "bytes"
    format "fmt"
    "github.com/plexmediamanager/service/log"
    "os/exec"
    "runtime"
    "strings"
)

type OSType int

const (
    OS_WINDOWS              OSType          =   1
    OS_LINUX                OSType          =   2
    OS_MACOS                OSType          =   3
)

var operatingSystemType     OSType

var operatingSystemMap      map[OSType]string

func init() {
    switch runtime.GOOS {
        case "windows":
            operatingSystemType = OS_WINDOWS
        case "linux":
            operatingSystemType = OS_LINUX
        case "darwin":
            operatingSystemType = OS_MACOS
        default:
            log.Panic("Unsupported OS detected:", runtime.GOOS)
    }
    operatingSystemMap = map[OSType]string {
        OS_WINDOWS:         "Windows",
        OS_LINUX:           "Linux",
        OS_MACOS:           "MacOS",
    }
}

type ServerInformation struct {
    Kernel                  *KernelInformation          `json:"kernel"`
    Processor               *ProcessorInformation       `json:"processor"`
    Memory                  *MemoryInformation          `json:"memory"`
    Network                 *NetworkInformation         `json:"network"`
    Uptime                  string                      `json:"uptime"`
}

func execute(commandName string,  arguments ...string) (string, error) {
    var errorOut bytes.Buffer
    command := exec.Command(commandName, arguments...)
    command.Stderr = &errorOut
    commandOutput, err := command.Output()
    outputToString := strings.TrimSpace(string(commandOutput))
    if err != nil {
        err = format.Errorf("%s: error=%q stderr=%s", commandName, err, string(errorOut.Bytes()))
    }
    return outputToString, err
}

// Convert bytes to human readable
func bytesToHumanReadable(bytes uint64) string {
    const unit = 1024
    if bytes < unit {
        return format.Sprintf("%d B", bytes)
    }
    div, exp := int64(unit), 0
    for n := bytes / unit; n >= unit; n /= unit {
        div *= unit
        exp++
    }
    return format.Sprintf("%.1f %ciB", float64(bytes) / float64(div), "KMGTPE"[exp])
}

func GetServerInformation() *ServerInformation {
    return &ServerInformation {
        Kernel:         getKernelInformation(),
        Processor:      getProcessorInformation(),
        Memory:         getMemoryInformation(),
        Network:        getNetworkInformation(),
        Uptime:         getSystemUptime(),
    }
}

func getSystemUptime() string {
    var uptime string
    switch operatingSystemType {
        case OS_WINDOWS:
            uptime = "not supported"
        case OS_LINUX:
            tmpTime, err := execute("uptime", "-p")
            if err == nil {
                uptime = tmpTime
            }
        case OS_MACOS:

    }
    return uptime
}