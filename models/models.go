package models

type Link struct {
	ID		string	`json:"id"`
	URL		*string	`json:"url"`
	Slug	*string	`json:"slug"`
}

func StrPt(s string) *string { 
	return &s 
}