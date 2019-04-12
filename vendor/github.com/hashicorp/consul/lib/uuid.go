package lib

import (
	"github.com/hashicorp/go-uuid"
)

// UUIDCheckFunc should determine whether the given UUID is actually
// unique and allowed to be used
type UUIDCheckFunc func(string) (bool, error)

func GenerateUUID(checkFn UUIDCheckFunc) (string, error) {
	for {
		id, err := uuid.GenerateUUID()
		if err != nil {
			return "", err
		}

		if checkFn == nil {
			return id, nil
		}

		if ok, err := checkFn(id); err != nil {
			return "", err
		} else if ok {
			return id, nil
		}
	}
}
