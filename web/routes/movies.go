package routes

import (
    "encoding/json"
    format "fmt"
    "github.com/plexmediamanager/micro-database/models"
    databaseService "github.com/plexmediamanager/micro-database/service"
    "github.com/plexmediamanager/micro-manager/web/routes/common"
    redisService "github.com/plexmediamanager/micro-redis/service"
    "github.com/plexmediamanager/service/errors"
    "net/http"
)

var (
    moviesModel                 models.Movie
    moviesList                  []*models.Movie

    moviesRedisKeyAll           string
    moviesRedisKeyDownloaded    string
)

type Movie struct {
    ID                          uint64                  `json:"id"`
    Title                       string                  `json:"title"`
    OriginalTitle               string                  `json:"original_title"`
    LocalTitle                  string                  `json:"local_title"`
    OriginalLanguage            string                  `json:"original_language"`
    Languages                   []*common.Language      `json:"languages"`
    Overview                    string                  `json:"overview"`
    Tagline                     string                  `json:"tagline"`
    Genres                      []*common.Genre          `json:"genres"`
    Homepage                    string                  `json:"homepage"`
    Runtime                     uint64                  `json:"runtime"`
    Status                      uint64                  `json:"status"`
    Adult                       bool                    `json:"adult"`
    ImdbId                      string                  `json:"imdb_id"`
    ReleaseDate                 string                  `json:"release_date"`
    //ProductionCompanies
    //ProductionCountries
    VoteAverage                 float64                 `json:"vote_average"`
    VoteCount                   uint64                  `json:"vote_count"`
    Popularity                  float64                 `json:"popularity"`
    Budget                      uint64                  `json:"budget"`
    Revenue                     uint64                  `json:"revenue"`
    Backdrop                    map[string]string       `json:"backdrop"`
    Poster                      map[string]string       `json:"poster"`
}

func InitializeMovieRouter() {
    moviesModel = models.Movie{}
    moviesRedisKeyAll = buildRedisCacheKey("models", moviesModel.TableName(), "all")
    moviesRedisKeyDownloaded = buildRedisCacheKey("models", moviesModel.TableName(), "downloaded")
    //moviesStoreRestoreAllToFromRedis()
    moviesStoreRestoreDownloadedToFromRedis()
}

func HandleMoviesList(writer http.ResponseWriter, request *http.Request) {
    response, err := databaseService.MovieServiceFindDownloaded(microClient())
    if err != nil {
        sendError(writer, &FailedResponse {
            Message:     "Failed to fetch list of movies",
            Data:        errors.ParseError(err),
        }, http.StatusOK)
    } else {
        moviesList := make([]*Movie, 0)
        for _, movie := range response {
            moviesList = append(moviesList, moviesBuildMovieInformation(movie))
        }
        sendResponse(writer, &SuccessfulResponse {
            Message:     "Successfully fetched list of movies",
            Data:        moviesList,
        }, http.StatusOK)
    }
}

// Store all movies to Redis
func moviesStoreRestoreAllToFromRedis() {
    has, err := redisService.RedisHas(microClient(), moviesRedisKeyAll)
    if err == nil {
        if !has {
            modelsList, err := databaseService.MovieServiceFindAll(microClient())
            if err != nil {
                format.Println(err)
            } else {
                moviesList = modelsList
                serialized, err := json.Marshal(modelsList)
                if err != nil {
                    format.Println(err)
                } else {
                    _, err = redisService.RedisSetValue(microClient(), moviesRedisKeyAll, string(serialized), 0)
                    if err != nil {
                        format.Println(err)
                    }
                }
            }
        } else {
            redisValue, err := redisService.RedisGetValue(microClient(), moviesRedisKeyAll)
            if err != nil {
                format.Println(err)
            } else {
                err := json.Unmarshal(redisValue.Response, &moviesList)
                if err != nil {
                    format.Println(err)
                }
            }
        }
    } else {
        format.Println(err)
    }
}

// Store/Restore downloaded movies to/from Redis
func moviesStoreRestoreDownloadedToFromRedis() {
    has, err := redisService.RedisHas(microClient(), moviesRedisKeyDownloaded)
    if err == nil {
        if !has {
            modelsList, err := databaseService.MovieServiceFindDownloaded(microClient())
            if err != nil {
                format.Println(err)
            } else {
                moviesList = modelsList
                serialized, err := json.Marshal(modelsList)
                if err != nil {
                    format.Println(err)
                } else {
                    _, err = redisService.RedisSetValue(microClient(), moviesRedisKeyDownloaded, string(serialized), 0)
                    if err != nil {
                        format.Println(err)
                    }
                }
            }
        } else {
            redisValue, err := redisService.RedisGetValue(microClient(), moviesRedisKeyDownloaded)
            if err != nil {
                format.Println(err)
            } else {
                err := json.Unmarshal(redisValue.Response, &moviesList)
                if err != nil {
                    format.Println(err)
                }
            }
        }
    } else {
        format.Println(err)
    }
}

func moviesBuildMovieInformation(item *models.Movie) *Movie {
    var itemGenres []uint64
    var itemLanguages []string
    _ = json.Unmarshal(item.Genres, &itemGenres)
    _ = json.Unmarshal(item.Languages, &itemLanguages)
    return &Movie {
        ID:                     item.ID,
        Title:                  item.Title,
        OriginalTitle:          item.OriginalTitle,
        LocalTitle:             item.LocalTitle,
        OriginalLanguage:       item.OriginalLanguage,
        Languages:              common.LoadLanguages(itemLanguages),
        Overview:               item.Overview,
        Tagline:                item.Tagline,
        Genres:                 common.LoadGenres(itemGenres),
        Homepage:               item.Homepage,
        Runtime:                item.Runtime,
        Status:                 item.Status,
        Adult:                  item.Adult,
        ImdbId:                 item.ImdbId,
        ReleaseDate:            item.ReleaseDate,
        //ProductionCompanies
        //ProductionCountries
        VoteAverage:            item.VoteAverage,
        VoteCount:              item.VoteCount,
        Popularity:             item.Popularity,
        Budget:                 item.Budget,
        Revenue:                item.Revenue,
        Backdrop:               common.BuildImagePath(item.Backdrop, "backdrop"),
        Poster:                 common.BuildImagePath(item.Poster, "poster"),
    }
}