package testutils

import (
	"io/ioutil"

	"github.com/Mind-Informatica-srl/restapi/pkg/delegate"
)

type SimpleObject struct {
	Nome    string `json:"nome" gorm:"primaryKey"`
	Cognome string `json:"cognome"`
}

var SimpleObjectDelegate = delegate.Delegate{
	ObjectCreator: func() interface{} {
		return &SimpleObject{}
	},
	ListCreator: func() interface{} {
		return &[]SimpleObject{}
	},
}

func SomeSimpleObject() ([]byte, error) {
	return ioutil.ReadFile("../../internal/testutils/testdata/simple_object.json")
}
