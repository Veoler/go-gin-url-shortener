package data

import (
	"github.com/Veoler/go-gin-url-shortener/models"
	"errors"
	"net/url"
)

var links = []models.Link {
	{
		ID: 1,
		Slug: models.StrPt("course"),
		URL: models.StrPt("https://example.com/super-long-course-url"),
	},
	{
		ID: 2,
		Slug: models.StrPt("blog"),
		URL: models.StrPt("https://example.com/blog"),
	},
}

var nextID = 3

func GetAll() ([]models.Link, error){
	return links, nil
}

func GetBySlug(slug string) (models.Link, error) {
	for _, r := range links {
		if *r.Slug == slug {
			return r, nil
		}
	}
	return models.Link{}, errors.New("Ссылка не найдена")
}

func Update(id int, input models.Link) (models.Link, error) {
	for i := range links {
		if links[i].ID == id {
			if input.Slug != nil {
				links[i].Slug = input.Slug
			}
			if input.URL != nil {
				links[i].URL = input.URL
			} 
			return links[i], nil
		}
	}
	return models.Link{}, errors.New("Ссылка не найдена")
}

func Add(l models.Link) models.Link {
	l.ID = nextID
	nextID++
	links = append(links, l)
	return l 
}

func DeleteID(id int) error {
	for i, d := range links {
		if d.ID == id {
			links = append(links[:i], links[i+1:]...)
			return nil
		}
	}
	return errors.New("ID не найден")
}

func DeleteSlug(slug string) error {
	for i, d := range links {
		if *d.Slug == slug {
			links = append(links[:i], links[i+1:]...)
			return nil
		}
	}
	return errors.New("Slug не найден")
}

func Redirect(slug string) (string, error) {
	for _, r := range links {
		if r.Slug != nil && *r.Slug == slug {
			return *r.URL, nil
		}
	}
	return "", errors.New("Slug не найден")
}

func IsSlugExist(slug string) bool {
	for _, r := range links {
		if r.Slug != nil && *r.Slug == slug {
			return true
		}
	}
	return false
}

func IsSlugExistByOther(slug string, currentID int) bool {
    for _, r := range links {
        if r.ID != currentID && r.Slug != nil && *r.Slug == slug {
            return true
        }
    }
    return false
}

func IsInvalidURL(u *string) bool {
    if u == nil || *u == "" { 
		return false 
	}
    _, err := url.ParseRequestURI(*u)
    return err != nil
}

// func isValidURL(raw string) bool {
//     parsed, err := url.ParseRequestURI(raw)
//     if err != nil {
//         return false
//     }

//     return true
// }
