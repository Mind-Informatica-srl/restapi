package models

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"

	"github.com/Mind-Informatica-srl/restapi/pkg/constants"
)

// https://stackoverflow.com/questions/42037562/golang-gorm-time-data-type-conversion

// OnlyTime time.Time per solo ora HH:MM
type OnlyTime time.Time

// UnmarshalJSON Implement Unmarshaler interface
func (j *OnlyTime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(constants.HourFormatStringHHMM, s)
	if err != nil {
		return err
	}
	*j = OnlyTime(t)

	fmt.Println(t)
	return nil
}

// MarshalJSON Implement Marshaler interface
func (j OnlyTime) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(j).Format(constants.HourFormatStringHHMM) + "\""), nil

}

// Format function for printing your date
func (j OnlyTime) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}

func NewOnlyTime(hour, min, sec int) OnlyTime {
	t := time.Date(0, time.January, 1, hour, min, sec, 0, time.UTC)
	return OnlyTime(t)
}

// Scan per OnlyTime
func (t *OnlyTime) Scan(value interface{}) error {
	switch v := value.(type) {
	case []byte:
		return t.UnmarshalText(string(v))
	case string:
		return t.UnmarshalText(v)
	case time.Time:
		*t = OnlyTime(v)
	case nil:
		*t = OnlyTime{}
	default:
		return fmt.Errorf("cannot sql.Scan() OnlyTime from: %#v", v)
	}
	return nil
}

// Value per OnlyTime
func (t OnlyTime) Value() (driver.Value, error) {
	return driver.Value(time.Time(t).Format(constants.HourFormatStringHHMMSS)), nil
}

// UnmarshalText per OnlyTime
func (t *OnlyTime) UnmarshalText(value string) error {
	dd, err := time.Parse(constants.HourFormatStringHHMMSS, value)
	if err != nil {
		return err
	}
	*t = OnlyTime(dd)
	return nil
}

// GormDataType per OnlyTime
// serve a gorm per sapere il tipo sql
func (OnlyTime) GormDataType() string {
	return "TIME without time zone"
}
