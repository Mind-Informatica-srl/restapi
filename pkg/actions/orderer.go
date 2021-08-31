package actions

import (
	"net/http"

	"gorm.io/gorm"
)

// Ordinate restituisce scope per ordinare una query
func Ordinate(r *http.Request) (paginateScope func(db *gorm.DB) *gorm.DB, sort string, direction string) {
	if params, ok := r.URL.Query()["sort"]; ok {
		sort = params[0]
	}
	if params, ok := r.URL.Query()["order"]; ok {
		direction = params[0]
	}

	paginateScope = func(db *gorm.DB) *gorm.DB {
		if sort != "" {
			return db.Order(ToSnakeCase(sort) + " " + direction)
		}
		return db
	}
	return
}
