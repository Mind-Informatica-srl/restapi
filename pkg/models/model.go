// PAckage models provide some utils stuff to manage models in restapi server
package models

// PKModel expose functions needed to manage the primary key in a model
type PKModel interface {
	// ExtractPK return the pk of the model
	ExtractPK() interface{}
	// SetPK set the pk for the model
	SetPK(pk interface{}) error
	// VerifyPK check the pk value
	VerifyPK(pk interface{}) (bool, error)
}
