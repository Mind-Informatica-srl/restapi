package models

import (
	"testing"

	"github.com/Mind-Informatica-srl/restapi/pkg/testutils"
)

func TestPkModel(t *testing.T) {
	var obj testutils.SimpleObjectWithId
	id := 1
	if err := obj.SetPK(id); err != nil {
		t.Error(err)
		t.Fail()
	}
	if obj.ID != id {
		t.Log("obj ID should be equal to id value")
		t.Fail()
	}
	if ok, err := obj.VerifyPK(id); err != nil {
		t.Log(err)
		t.Fail()
	} else if !ok {
		t.Log("obj ID should be verified")
		t.Fail()
	}
}
