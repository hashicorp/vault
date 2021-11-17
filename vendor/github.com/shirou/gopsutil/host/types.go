package host

import (
	"fmt"
)

type Warnings struct {
	List []error
}

func (w *Warnings) Add(err error) {
	w.List = append(w.List, err)
}

func (w *Warnings) Reference() error {
	if len(w.List) > 0 {
		return w
	} else {
		return nil
	}
}

func (w *Warnings) Error() string {
	return fmt.Sprintf("Number of warnings: %v", len(w.List))
}
