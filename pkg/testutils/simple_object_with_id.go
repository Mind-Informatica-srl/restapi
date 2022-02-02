package testutils

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type SimpleObjectWithId struct {
	ID      int
	Nome    string
	Cognome string
}

// SetPK set the pk for the model
func (o *SimpleObjectWithId) SetPK(pk interface{}) error {
	o.ID = pk.(int)
	return nil
}

// VerifyPK check the pk value
func (o *SimpleObjectWithId) VerifyPK(pk interface{}) (bool, error) {
	return o.ID == pk.(int), nil
}

type SimpleObjectWithIdDelegate struct {
	DB *gorm.DB
}

func (d SimpleObjectWithIdDelegate) ProvideDB() *gorm.DB {
	return d.DB
}

func (d SimpleObjectWithIdDelegate) CreateObject(r *http.Request) (interface{}, error) {
	return &SimpleObjectWithId{}, nil
}

func (d SimpleObjectWithIdDelegate) ExtractPK(r *http.Request) (map[string]interface{}, error) {
	vars := mux.Vars(r)
	if pk, err := strconv.Atoi(vars["id"]); err != nil {
		return nil, err
	} else {
		res := make(map[string]interface{})
		res["id"] = pk
		return res, nil
	}
}

func (d SimpleObjectWithIdDelegate) VerifyPK(element interface{}, pk map[string]interface{}) (bool, error) {
	e := element.(*SimpleObjectWithId)
	id := pk["id"].(int)
	if e.ID != id {
		return false, nil
	}
	return true, nil
}

func (d SimpleObjectWithIdDelegate) AssignPK(element interface{}, pk interface{}) error {
	e := element.(*SimpleObjectWithId)
	id := pk.(int)
	e.ID = id
	return nil
}

func SomeSimpleObjectWithId() ([]byte, error) {
	return ioutil.ReadFile("../../internal/testutils/testdata/simple_object_with_id.json")
}
