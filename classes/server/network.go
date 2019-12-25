package server

import (
    "encoding/json"
    "github.com/plexmediamanager/service/helpers"
    "io/ioutil"
    "net"
    "net/http"
    "strings"
)

type backendFrontendNetwork struct {
    RemoteIP            string                  `json:"remote_ip"`
    LocalIP             string                  `json:"local_ip"`
    Domain              string                  `json:"domain"`
}

type nameServers struct {
    Domain              string                  `json:"domain"`
    IP                  string                  `json:"ip"`
}

type NetworkInformation struct {
    Backend             backendFrontendNetwork  `json:"backend"`
    Frontend            backendFrontendNetwork  `json:"frontend"`
    NameServers         []*nameServers           `json:"nameservers"`
}

func getNetworkInformation() *NetworkInformation {
    remoteAddress := getRemoteIP()
    return &NetworkInformation {
        Backend:        backendFrontendNetwork {
            RemoteIP:   remoteAddress,
            LocalIP:    getLocalIP(),
            Domain:     helpers.GetEnvironmentVariableAsString("BACKEND_DOMAIN", ""),
        },
        Frontend:       backendFrontendNetwork {
            RemoteIP:   remoteAddress,
            Domain:     helpers.GetEnvironmentVariableAsString("FRONTEND_DOMAIN", ""),
        },
        NameServers:    resolveNameServersForAddress(remoteAddress),
    }
}

// Get local IP address
func getLocalIP() string {
    conn, err := net.Dial("udp", "8.8.8.8:80")
    if err != nil {
        return ""
    }
    defer conn.Close()
    localAddr := conn.LocalAddr().String()
    idx := strings.LastIndex(localAddr, ":")
    return localAddr[0:idx]
}

// Get Remote IP address
func getRemoteIP() string {
    var httpClient http.Client

    response, err := httpClient.Get("https://ifconfig.co/json")
    if err != nil {
        return ""
    }

    result, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return ""
    }

    type tempStructure struct {
        IP          string      `json:"ip"`
    }
    var temp tempStructure
    err = json.Unmarshal(result, &temp)

    if err != nil {
        return ""
    }
    return temp.IP
}

// Resolve nameservers for address
func resolveNameServersForAddress(remoteIP string) []*nameServers {
    ns := make([]*nameServers, 0)
    serverNames, _ := net.LookupAddr(remoteIP)
    for _, serverName := range serverNames {
        name := serverName[:len(serverName) - 1]
        nServers, _ := net.LookupNS(name)
        for _, nameServer := range nServers {
            nsHost := nameServer.Host[:len(nameServer.Host) - 1]
            nsAddresses, _ := net.LookupIP(nsHost)
            for _, address := range nsAddresses {
                ns = append(ns, &nameServers {
                    Domain: nsHost,
                    IP:     address.String(),
                })
            }
        }
    }
    return ns
}