package parseabletime

import (
	"time"
)

const (
	dateLayout = "2006-01-02T15:04:05"
)

type ParseableTime time.Time

func (p *ParseableTime) UnmarshalJSON(b []byte) error {
	t, err := time.Parse(`"`+dateLayout+`"`, string(b))
	if err != nil {
		return err
	}

	*p = ParseableTime(t)

	return nil
}
