package restapi

import "encoding/json"

type Pager struct {
	TotalCount int64
	Items      func() interface{}
}

func (p Pager) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		TotalCount int64
		Items      interface{}
	}{
		TotalCount: p.TotalCount,
		Items:      p.Items(),
	})
}
