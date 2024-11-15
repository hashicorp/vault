package pgx

import (
	"database/sql/driver"

	"github.com/jackc/pgtype"
)

func convertDriverValuers(args []interface{}) ([]interface{}, error) {
	for i, arg := range args {
		switch arg := arg.(type) {
		case pgtype.BinaryEncoder:
		case pgtype.TextEncoder:
		case driver.Valuer:
			v, err := callValuerValue(arg)
			if err != nil {
				return nil, err
			}
			args[i] = v
		}
	}
	return args, nil
}
