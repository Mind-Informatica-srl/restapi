// Package models raggruppa model utili
package models

import (
	"database/sql/driver"
	"encoding/json"
)

type JSONB map[string]interface{}

func (j JSONB) Value() (driver.Value, error) {
	valueString, err := json.Marshal(j)
	return string(valueString), err
}

func (j *JSONB) Scan(value interface{}) error {
	if err := json.Unmarshal(value.([]byte), &j); err != nil {
		return err
	}
	return nil
}

// GormDataType per JSONB
// serve a gorm per sapere il tipo sql
func (JSONB) GormDataType() string {
	return "jsonb"
}

// ToJSONB trasforma un'interface in un jsonb o restituisce un errore
func ToJSONB(value interface{}) (obj JSONB, err error) {
	var buf []byte
	if buf, err = json.Marshal(value); err != nil {
		return
	}
	err = json.Unmarshal(buf, &obj)
	return
}
