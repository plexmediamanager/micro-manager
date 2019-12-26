package common

import databaseService "github.com/plexmediamanager/micro-database/service"

type Genre struct {
    ID                          uint64              `json:"id"`
    Name                        string              `json:"name"`
}

// Load genres
func LoadGenres(genres []uint64) []*Genre {
    genresList := make([]*Genre, 0)
    if len(genres) == 0 {
        return genresList
    }
    result, err := databaseService.GenreServiceFindManyByID(microClient(), genres...)
    if err != nil {
        return nil
    } else {
        for _, item := range result {
            genresList = append(genresList, &Genre {
                ID:     item.ID,
                Name:   item.Name,
            })
        }
    }
    return genresList
}
