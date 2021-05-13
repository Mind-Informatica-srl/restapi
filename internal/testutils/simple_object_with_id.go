package testutils

import (
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/Mind-Informatica-srl/restapi/pkg/delegate"
	"github.com/gorilla/mux"
)

type SimpleObjectWithId struct {
	ID      int
	Nome    string
	Cognome string
}

var SimpleObjectWithIdDelegate = delegate.Delegate{
	ObjectCreator: func() interface{} {
		return &SimpleObjectWithId{}
	},
	PKExtractor: func(r *http.Request) (interface{}, error) {
		vars := mux.Vars(r)
		if pk, err := strconv.Atoi(vars["id"]); err != nil {
			return nil, err
		} else {
			return pk, nil
		}
	},
	PKVerificator: func(element interface{}, pk interface{}) (bool, error) {
		e := element.(*SimpleObjectWithId)
		id := pk.(int)
		if e.ID != id {
			return false, nil
		}
		return true, nil
	},
	PKAssigner: func(element interface{}, pk interface{}) {
		e := element.(*SimpleObjectWithId)
		id := pk.(int)
		e.ID = id
	},
}

func SomeSimpleObjectWithId() ([]byte, error) {
	return ioutil.ReadFile("../../internal/testutils/testdata/simple_object_with_id.json")
}
