package server

import (
    "github.com/plexmediamanager/service/log"
    "strconv"
    "strings"
)

type memoryFormatted struct {
    Exact               uint64              `json:"exact"`
    Nice                string              `json:"nice"`
}

type MemoryInformation struct {
    Total               memoryFormatted     `json:"total"`
    Free                memoryFormatted     `json:"free"`
    Used                memoryFormatted     `json:"used"`
    Available           memoryFormatted     `json:"available"`
    Cached              memoryFormatted     `json:"cached"`
}

func getWindowsTotalMemorySize() (uint64, uint64) {
    var err error
    result, err := execute("wmic", "os", "get", "FreePhysicalMemory,", "TotalVisibleMemorySize")
    if err != nil {
        log.Panic(err)
    }
    resultToSlice := strings.Split(strings.Split(result, "\r\r\n")[1], "  ")
    memoryDetails := make([]uint64, 0)
    for _, item := range resultToSlice {
        if len(strings.TrimSpace(item)) > 0 {
            value, err := strconv.ParseInt(strings.TrimSpace(item), 10, 64)
            if err == nil {
                memoryDetails = append(memoryDetails, uint64(value) * 1000)
            }
        }
    }
    return memoryDetails[0], memoryDetails[1]
}

func getLinuxMemoryInformation(totalMemory *uint64, freeMemory *uint64, availableMemory *uint64, cachedMemory *uint64) {
    var err error
    result, err := execute("sh",
        "-c",
        "cat /proc/meminfo",
    )
    if err != nil {
        log.Panic(err)
    }
    partedResult := strings.Split(result, "\n")
    for _, part := range partedResult {
        total := extractValueFromString(part, "MemTotal")
        if len(total) > 0 {
            value, err := strconv.ParseUint(total, 10, 64)
            if err == nil {
                *totalMemory = value * 1024
            }
        }

        free := extractValueFromString(part, "MemFree")
        if len(free) > 0 {
            value, err := strconv.ParseUint(free, 10, 64)
            if err == nil {
                *freeMemory = value * 1024
            }
        }

        available := extractValueFromString(part, "MemAvailable")
        if len(available) > 0 {
            value, err := strconv.ParseUint(available, 10, 64)
            if err == nil {
                *availableMemory = value * 1024
            }
        }

        cached := extractValueFromString(part, "Cached")
        if len(cached) > 0 {
            value, err := strconv.ParseUint(cached, 10, 64)
            if err == nil {
                *cachedMemory = value * 1024
            }
        }
    }
}

func getMemoryInformation() *MemoryInformation {
    var totalMemory uint64
    var freeMemory uint64
    var usedMemory uint64
    var availableMemory uint64
    var cachedMemory uint64
    switch operatingSystemType {
        case OS_WINDOWS:
            freeMemory, totalMemory = getWindowsTotalMemorySize()
            usedMemory = totalMemory - freeMemory
            cachedMemory = freeMemory - usedMemory
            availableMemory = freeMemory
        case OS_LINUX:
            getLinuxMemoryInformation(&totalMemory, &freeMemory, &availableMemory, &cachedMemory)
            usedMemory = totalMemory - availableMemory
        case OS_MACOS:

    }
    return &MemoryInformation {
        Total:      memoryFormatted {
            Exact:  totalMemory,
            Nice:   bytesToHumanReadable(totalMemory),
        },
        Free:       memoryFormatted {
            Exact:  freeMemory,
            Nice:   bytesToHumanReadable(freeMemory),
        },
        Used:       memoryFormatted {
            Exact:  usedMemory,
            Nice:   bytesToHumanReadable(usedMemory),
        },
        Available:  memoryFormatted {
            Exact:  availableMemory,
            Nice:   bytesToHumanReadable(availableMemory),
        },
        Cached:     memoryFormatted {
            Exact:  cachedMemory,
            Nice:   bytesToHumanReadable(cachedMemory),
        },
    }
}

func extractValueFromString(row string, lookFor string) string {
    if strings.HasPrefix(row, lookFor) {
        return strings.TrimSpace(strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(row, "kB", ""), ":", ""), lookFor, ""))
    }
    return ""
}