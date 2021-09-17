package testutils

import (
	"io/ioutil"

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

func (d SimpleObjectDelegate) CreateObject() interface{} {
	return &SimpleObject{}
}

func (d SimpleObjectDelegate) CreateList() interface{} {
	return &[]SimpleObject{}
}

func SomeSimpleObject() ([]byte, error) {
	return ioutil.ReadFile("../../internal/testutils/testdata/simple_object.json")
}
