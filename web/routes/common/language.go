package common

import databaseService "github.com/plexmediamanager/micro-database/service"

type Language struct {
    ISO                 string          `json:"iso"`
    Name                string          `json:"name"`
}

func LoadLanguages(languages []string) []*Language {
    languagesList := make([]*Language, 0)
    if len(languages) == 0 {
        return languagesList
    }
    result, err := databaseService.LanguageServiceFindManyByISO(microClient(), languages...)
    if err != nil {
        return nil
    } else {
        for _, language := range result {
            var languageName string
            if len(language.Name) == 0 {
                languageName = language.EnglishName
            } else {
                languageName = language.Name
            }
            languagesList = append(languagesList, &Language {
                ISO:    language.Iso,
                Name:   languageName,
            })
        }
    }
    return languagesList
}