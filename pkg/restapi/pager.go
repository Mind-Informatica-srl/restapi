package restapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type Pager struct {
	TotalCount int64
	Items      func() interface{}
}

func (p Pager) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TotalCount int64       `json:"totalCount"`
		Items      interface{} `json:"items"`
	}{
		TotalCount: p.TotalCount,
		Items:      p.Items(),
	})
}

func Paginate(r *http.Request) (paginateScope func(db *gorm.DB) *gorm.DB, page int, pageSize int) {
	var pageString, pageSizeString string
	if params, ok := r.URL.Query()["page"]; ok {
		pageString = params[0]
	}
	if params, ok := r.URL.Query()["pageSize"]; ok {
		pageSizeString = params[0]
	}
	if pageString != "" && pageSizeString != "" {
		pageSize, _ = strconv.Atoi(pageSizeString)
		page, _ = strconv.Atoi(pageString)
	}
	paginateScope = func(db *gorm.DB) *gorm.DB {
		if page > 0 && pageSize > 0 {
			offset := (page - 1) * pageSize
			return db.Offset(offset).Limit(pageSize)
		}
		return db
	}
	return
}
