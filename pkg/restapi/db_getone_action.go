package restapi

import (
	"encoding/json"
	"net/http"
)

type DBGetOneAction DatabaseAction

func (action *DBGetOneAction) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	db := action.DBProvider()
	id, err := action.PKExtractor(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	element := action.ObjectCreator()
	if err := db.First(element, id).Error; err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(element); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
