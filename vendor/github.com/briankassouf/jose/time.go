package jose

import "time"

// Now returns the current time in UTC.
func Now() time.Time { return time.Now().UTC() }
