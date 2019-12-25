package server

import (
    format "fmt"
    "github.com/plexmediamanager/service/log"
    "runtime"
    "strconv"
    "strings"
)

type ProcessorInformation struct {
    Vendor                      string      `json:"vendor"`
    Model                       string      `json:"model"`
    Cores                       int         `json:"cores"`
    Threads                     int         `json:"threads"`
    Frequency                   int64       `json:"frequency"`
}

// Get processor information for Windows
func getWindowsProcessorInformation(vendor *string, model *string, cores *int, frequency *int64) {
    var err error
    result, err := execute("wmic", "cpu", "get", "name,", "numberofcores,", "maxclockspeed")
    if err != nil {
        log.Panic(err)
    }
    data := strings.Split(result, "\n")[1]
    details := make([]string, 0)
    for _, item := range strings.Split(data, "  ") {
        if len(strings.TrimSpace(item)) > 0 {
            details = append(details, strings.TrimSpace(item))
        }
    }

    *frequency, err = strconv.ParseInt(details[0], 10, 64)
    if err != nil {
        *frequency = 800
    }
    *cores, err = strconv.Atoi(details[2])
    if err != nil {
        *cores = runtime.NumCPU() / 2
    }
    processorName := details[1]
    var processorNameWithoutVendor string
    if strings.Contains(processorName, "Intel") {
        *vendor = "Intel"
        processorNameWithoutVendor = strings.TrimSpace(strings.ReplaceAll(processorName, "Intel(R) Core(TM)", ""))
        indexOfAt := strings.Index(processorNameWithoutVendor,"@")
        *model = strings.TrimSpace(strings.ReplaceAll(processorNameWithoutVendor[:indexOfAt], "CPU", ""))
    } else {
        *vendor = "AMD"
        processorNameWithoutVendor = strings.TrimSpace(strings.ReplaceAll(processorName, "AMD", ""))
        *model = strings.TrimSpace(strings.ReplaceAll(processorNameWithoutVendor, format.Sprintf("%d-Core Processor", *cores), ""))
    }
}

// Get processor information for Linux
func getLinuxProcessorInformation(vendor *string, model *string, cores *int, frequency *int64) {
    var err error

    coreCountResult, err := execute("sh",
        "-c",
        "cat /proc/cpuinfo | grep 'cpu cores' | uniq | cut -f2 -d ':'",
    )
    if err != nil {
        log.Panic(err)
    }
    *cores, err = strconv.Atoi(strings.TrimSpace(coreCountResult))
    if err != nil {
        *cores = runtime.NumCPU() / 2
    }

    vendorResult, err := execute("sh",
        "-c",
        "cat /proc/cpuinfo | grep 'model name' | uniq | cut -f2 -d ':'",
    )
    if err != nil {
        log.Panic(err)
    }
    var processorNameWithoutVendor string
    if strings.Contains(vendorResult, "Intel") {
        *vendor = "Intel"
        processorNameWithoutVendor = strings.TrimSpace(strings.ReplaceAll(vendorResult, "Intel(R) Core(TM)", ""))
        indexOfAt := strings.Index(processorNameWithoutVendor,"@")
        *model = strings.TrimSpace(strings.ReplaceAll(processorNameWithoutVendor[:indexOfAt], "CPU", ""))
    } else {
        *vendor = "AMD"
        processorNameWithoutVendor = strings.TrimSpace(strings.ReplaceAll(vendorResult, "AMD", ""))
        *model = strings.TrimSpace(strings.ReplaceAll(processorNameWithoutVendor, format.Sprintf("%d-Core Processor", *cores), ""))
    }
    cpuFrequencyResult, err := execute("sh",
       "-c",
       "cat /proc/cpuinfo | grep 'cpu MHz' | uniq | cut -f2 -d ':'",
    )
    if err != nil {
       log.Panic(err)
    }
    for _, freq := range strings.Split(cpuFrequencyResult, "\n") {
        value, err := strconv.ParseFloat(freq, 64)
        if err == nil {
            *frequency = int64(value)
            break
        }
    }
}

func getProcessorInformation() *ProcessorInformation {
    var vendor string
    var model string
    var cores int
    threads := runtime.NumCPU()
    var frequency int64

    switch operatingSystemType {
        case OS_WINDOWS:
            getWindowsProcessorInformation(&vendor, &model, &cores, &frequency)
        case OS_LINUX:
            getLinuxProcessorInformation(&vendor, &model, &cores, &frequency)
        case OS_MACOS:

    }

    return &ProcessorInformation {
        Vendor:     vendor,
        Model:      model,
        Cores:      cores,
        Threads:    threads,
        Frequency:  frequency,
    }
}