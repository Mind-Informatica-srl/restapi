package testutils

import (
	"io/ioutil"
	"net/http"

	"gorm.io/gorm"
)

type SimpleObject struct {
	Nome    string `json:"nome" gorm:"primaryKey"`
	Cognome string `json:"cognome"`
}

type SimpleObjectDelegate struct {
	DB *gorm.DB
}

func (d SimpleObjectDelegate) ProvideDB() *gorm.DB {
	return d.DB
}

func (d SimpleObjectDelegate) CreateObject(r *http.Request) (interface{}, error) {
	return &SimpleObject{}, nil
}

func (d SimpleObjectDelegate) CreateList(r *http.Request) (interface{}, error) {
	return &[]SimpleObject{}, nil
}

func (d SimpleObjectDelegate) CSVFileName(r *http.Request) (string, error) {
	return "export.csv", nil
}

func SomeSimpleObject() ([]byte, error) {
	return ioutil.ReadFile("../../internal/testutils/testdata/simple_object.json")
}
