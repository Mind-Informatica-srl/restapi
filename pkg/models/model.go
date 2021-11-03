// PAckage models provide some utils stuff to manage models in restapi server
package models

import (
	"net/http"

	"gorm.io/gorm"
)

// PKModel expose functions needed to manage the primary key in a model
type PKModel interface {
	// SetPK set the pk for the model
	SetPK(pk interface{}) error
	// VerifyPK check the pk value
	VerifyPK(pk interface{}) (bool, error)
}

// BaseDelegate implements the DBDelegate interfaces for a model implementig the PKModel interface
type BaseDelegate struct {
	DBProvider    func() *gorm.DB
	ObjectCreator func() PKModel
	ListCreator   func() interface{}
	PKExtractor   func(r *http.Request) (interface{}, error)
	PKUrlPart     *string
}

func (d BaseDelegate) ExtractPK(r *http.Request) (interface{}, error) {
	return d.PKExtractor(r)
}

func (d BaseDelegate) ProvideDB() *gorm.DB {
	return d.DBProvider()
}

func (d BaseDelegate) AssignPK(element interface{}, pk interface{}) error {
	e := element.(PKModel)
	return e.SetPK(pk)
}

func (d BaseDelegate) CreateObject() interface{} {
	return d.ObjectCreator()
}

func (d BaseDelegate) CreateList() interface{} {
	return d.ListCreator()
}

func (d BaseDelegate) VerifyPK(element interface{}, pk interface{}) (bool, error) {
	return element.(PKModel).VerifyPK(pk)
}

func (d *BaseDelegate) PKUrl() string {
	if d.PKUrlPart != nil {
		return *d.PKUrlPart
	}
	return "/{id}"
}

func NewBaseDelegateWithPKUrl(dbProvider func() *gorm.DB, objectCreator func() PKModel, listCreator func() interface{}, pkExtractor func(r *http.Request) (interface{}, error), pkUrl *string) BaseDelegate {
	return BaseDelegate{
		DBProvider:    dbProvider,
		ObjectCreator: objectCreator,
		ListCreator:   listCreator,
		PKExtractor:   pkExtractor,
		PKUrlPart:     pkUrl,
	}
}

func NewBaseDelegate(dbProvider func() *gorm.DB, objectCreator func() PKModel, listCreator func() interface{}, pkExtractor func(r *http.Request) (interface{}, error)) BaseDelegate {
	return NewBaseDelegateWithPKUrl(dbProvider, objectCreator, listCreator, pkExtractor, nil)
}
