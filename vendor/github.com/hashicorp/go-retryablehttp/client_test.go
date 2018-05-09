package retryablehttp

import (
	"fmt"
	"testing"
	"time"
)

func TestLinearJitterBackoff(t *testing.T) {
	min := time.Duration(1000) * time.Second
	max := time.Duration(2000) * time.Second
	for i := 0; i < 20; i++ {
		result := LinearJitterBackoff(min, max, i, nil).Seconds()
		if result < min.Seconds() {
			t.Fatal(fmt.Sprintf("result too low, got %f", result))
		}
		if result > max.Seconds() {
			t.Fatal(fmt.Sprintf("result too high, got %f", result))
		}
	}
}