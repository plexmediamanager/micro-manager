package common

import (
    format "fmt"
    "github.com/plexmediamanager/micro-tmdb/tmdb"
    "github.com/plexmediamanager/service/ctx"
    "strings"
)

// Build Image Path
func BuildImagePath(image string, mediaType string) map[string]string {
    paths := make(map[string]string, 0)
    tmdbConfiguration := ctx.Value("tmdbConfiguration").(*tmdb.APIConfiguration)
    baseUrl := strings.TrimSuffix(tmdbConfiguration.Images.SecureBaseURL, "/")

    switch mediaType {
        case "company":
            fallthrough
        case "network":
            for _, size := range tmdbConfiguration.Images.LogoSizes {
                paths[size] = format.Sprintf("%s/%s/%s", baseUrl, size, image)
            }
        case "creator":
            for _, size := range tmdbConfiguration.Images.ProfileSizes {
                paths[size] = format.Sprintf("%s/%s/%s", baseUrl, size, image)
            }
        case "backdrop":
            for _, size := range tmdbConfiguration.Images.BackdropSizes {
                paths[size] = format.Sprintf("%s/%s/%s", baseUrl, size, image)
            }
        case "poster":
            for _, size := range tmdbConfiguration.Images.PosterSizes {
                paths[size] = format.Sprintf("%s/%s/%s", baseUrl, size, image)
            }
    }

    return paths
}
