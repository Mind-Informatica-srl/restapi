package actions

import (
	"errors"
	"testing"
)

func TestActionErrorError(t *testing.T) {
	cases := []struct {
		name string
		err  ActionError
		want string
	}{
		{"nil err with status", ActionError{Err: nil, Status: 400}, "Bad Request"},
		{"zero value", ActionError{}, "Internal Server Error"},
		{"nil err with non standard status", ActionError{Status: 599}, "Internal Server Error"},
		{"nil nested ActionError", ActionError{Err: (*ActionError)(nil), Status: 404}, "Not Found"},
		{"wrapped error", ActionError{Err: errors.New("x"), Status: 400}, "x"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := c.err.Error(); got != c.want {
				t.Errorf("Error() = %q, want %q", got, c.want)
			}
		})
	}
}
