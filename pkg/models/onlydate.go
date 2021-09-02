package models

import (
	"strings"
	"time"

	"github.com/Mind-Informatica-srl/restapi/pkg/constants"
)

type OnlyDate time.Time

// Implement Marshaler and Unmarshaler interface
func (j *OnlyDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse(constants.DateFormatStringYYYYMMDD, s)
	if err != nil {
		// _, ok := err.(*time.ParseError)
		// if !ok {
		// 	return err
		// }
		if strings.Contains(s, "T") {
			s = strings.Split(s, "T")[0]
			t, err = time.Parse(constants.DateFormatStringYYYYMMDD, s)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	*j = OnlyDate(t)
	return nil
}

func (j OnlyDate) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Time(j).Format(constants.DateFormatStringYYYYMMDD) + "\""), nil

}

// Maybe a Format function for printing your date
func (j OnlyDate) Format(s string) string {
	t := time.Time(j)
	return t.Format(s)
}

// GormDataType per OnlyDate
// serve a gorm per sapere il tipo sql
func (OnlyDate) GormDataType() string {
	return "date"
}
