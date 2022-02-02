package models

import (
	"testing"

	"github.com/Mind-Informatica-srl/restapi/pkg/testutils"
)

func TestPkModel(t *testing.T) {
	var obj testutils.SimpleObjectWithId
	m := make(map[string]interface{})
	m["id"] = 1
	if err := obj.SetPK(m); err != nil {
		t.Error(err)
		t.Fail()
	}
	if obj.ID != m["id"] {
		t.Log("obj ID should be equal to id value")
		t.Fail()
	}
	if ok, err := obj.VerifyPK(m); err != nil {
		t.Log(err)
		t.Fail()
	} else if !ok {
		t.Log("obj ID should be verified")
		t.Fail()
	}
}
