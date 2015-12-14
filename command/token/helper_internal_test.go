package token

import (
	"testing"
)

// TestCommand re-uses the existing Test function to ensure proper behavior of
// the internal token helper
func TestCommand(t *testing.T) {
	Test(t, &InternalTokenHelper{})
}
