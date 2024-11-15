package internal

import "os"

// Get first non-empty value
func GetDefaultString(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}

	return ""
}

// set back the memoried enviroment variables
type Rollback func()

func Memory(keys ...string) Rollback {
	// remenber enviroment variables
	m := make(map[string]string)
	for _, key := range keys {
		m[key] = os.Getenv(key)
	}

	return func() {
		for _, key := range keys {
			os.Setenv(key, m[key])
		}
	}
}
