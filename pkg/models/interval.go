package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgtype"
)

// PGInterval tipo per comunicare con postgres per colonne
// di tipo interval
type PGInterval pgtype.Interval

// Interval struct scambiata con il client
// per intervalli
type Interval struct {
	Anni    int
	Mesi    int
	Giorni  int
	Ore     int
	Minuti  int
	Secondi int
}

// UnmarshalJSON Implement Unmarshaler interface
func (i *PGInterval) UnmarshalJSON(b []byte) error {
	var val Interval
	if err := json.Unmarshal(b, &val); err != nil {
		return err
	}
	if err := i.ConvertToPGInterval(val); err != nil {
		return err
	}
	fmt.Println(i)
	return nil
}

// MarshalJSON Implement Marshaler interface
func (j PGInterval) MarshalJSON() ([]byte, error) {
	res, err := j.ConvertFromPGInterval()
	if err != nil {
		return nil, err
	}
	return json.Marshal(res)
}

// Scan implements the database/sql Scanner interface.
func (dst *PGInterval) Scan(src interface{}) error {
	interval := pgtype.Interval(*dst)
	fmt.Println(interval)
	if err := interval.Scan(src); err != nil {
		return err
	}
	*dst = PGInterval(interval)
	return nil
}

// Value implements the database/sql/driver Valuer interface.
func (src PGInterval) Value() (driver.Value, error) {
	tmp, err := pgtype.Interval(src).Value()
	fmt.Println(tmp)
	return tmp, err
}

const (
	monthsPerYear         = 12
	hoursPerDay           = 24
	microsecondsPerSecond = 1000000
	microsecondsPerMinute = 60 * microsecondsPerSecond
	microsecondsPerHour   = 60 * microsecondsPerMinute
	microsecondsPerDay    = 24 * microsecondsPerHour
	microsecondsPerMonth  = 30 * microsecondsPerDay
)

func (i *PGInterval) ConvertToPGInterval(value Interval) (err error) {
	months := (int32(value.Anni) * monthsPerYear)
	months += int32(value.Mesi)
	days := int32(value.Giorni)
	microseconds := int64(value.Ore) * microsecondsPerHour
	microseconds += int64(value.Minuti) * microsecondsPerMinute
	microseconds += int64(value.Secondi) * microsecondsPerSecond

	var status pgtype.Status
	if days == 0 && microseconds == 0 && months == 0 {
		status = pgtype.Null
	} else {
		status = pgtype.Present
	}
	res := PGInterval{
		Microseconds: microseconds,
		Days:         days,
		Months:       months,
		Status:       status,
	}
	*i = PGInterval(res)
	return nil
}

func (j PGInterval) ConvertFromPGInterval() (res Interval, err error) {
	var anni, mesi, giorni, ore, minuti, secondi int

	ore = (int(j.Microseconds) / microsecondsPerHour)
	minuti = int(j.Microseconds%microsecondsPerHour) / microsecondsPerMinute
	secondi = int(j.Microseconds%microsecondsPerMinute) / microsecondsPerSecond
	// microseconds := j.Microseconds % microsecondsPerSecond

	giorni = int(j.Days) + ((int(j.Microseconds) / microsecondsPerHour) / microsecondsPerDay)

	anni = int(j.Months) / monthsPerYear
	mesi = int(j.Months) % monthsPerYear

	res = Interval{
		Anni:    anni,
		Mesi:    mesi,
		Giorni:  giorni,
		Ore:     ore,
		Minuti:  minuti,
		Secondi: secondi,
	}
	return
}
