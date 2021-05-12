package restapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"gorm.io/gorm"
)

type DBGetAllAction DatabaseAction

func (action *DBGetAllAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db := action.DBProvider()
	db, err := QueryFilter(db, r.URL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	list := action.ObjectCreator()
	var page, pageSize string
	var count int64
	if params, ok := r.URL.Query()["page"]; ok {
		page = params[0]
	}
	if params, ok := r.URL.Query()["pageSize"]; ok {
		pageSize = params[0]
	}
	//se è richiesta la paginazione
	if page != "" && pageSize != "" {
		var limit, offset int
		var err error
		if limit, err = strconv.Atoi(pageSize); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if offset, err = strconv.Atoi(page); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		offset = limit * (offset - 1)
		element := action.ObjectCreator()
		if err = db.Model(element).Count(&count).Error; err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		db = db.Limit(limit).Offset(offset)
	}
	if err := db.Find(list).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if page != "" && pageSize != "" {
		res := Pager{
			TotalCount: count,
			Items: func() interface{} {
				return list
			},
		}
		if err := json.NewEncoder(w).Encode(res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		if err := json.NewEncoder(w).Encode(list); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
