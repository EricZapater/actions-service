package models

import "github.com/google/uuid"

type Area struct {
	Id          uuid.UUID `json:"id"`
	Description string    `json:"description"`
	SiteId uuid.UUID `json:"siteId"`
}