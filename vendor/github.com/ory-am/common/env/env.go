// A very handy library which adds defaults to os.GetEnv()
package env

import "os"

// Getenv retrieves the value of the environment variable named by the key. It returns the value, which will return the fallback if the variable is not present.
func Getenv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
