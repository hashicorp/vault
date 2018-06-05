package lib

import (
	"math"
	"time"

	"github.com/hashicorp/serf/coordinate"
)

// ComputeDistance returns the distance between the two network coordinates in
// seconds. If either of the coordinates is nil then this will return positive
// infinity.
func ComputeDistance(a *coordinate.Coordinate, b *coordinate.Coordinate) float64 {
	if a == nil || b == nil {
		return math.Inf(1.0)
	}

	return a.DistanceTo(b).Seconds()
}

// CoordinateSet holds all the coordinates for a given node, indexed by network
// segment name.
type CoordinateSet map[string]*coordinate.Coordinate

// Intersect tries to return a pair of coordinates which are compatible with the
// current set and a given set. We employ some special knowledge about network
// segments to avoid doing a full intersection, since this is in several hot
// paths. This might return nil for either coordinate in the output pair if an
// intersection cannot be found. The ComputeDistance function above is designed
// to deal with that.
func (cs CoordinateSet) Intersect(other CoordinateSet) (*coordinate.Coordinate, *coordinate.Coordinate) {
	// Use the empty segment by default.
	segment := ""

	// If we have a single segment, then let our segment take priority since
	// we are possibly a client. Any node with more than one segment can only
	// be a server, which means it should be in all segments.
	if len(cs) == 1 {
		for s := range cs {
			segment = s
		}
	}

	// Likewise for the other set.
	if len(other) == 1 {
		for s := range other {
			segment = s
		}
	}

	return cs[segment], other[segment]
}

// GenerateCoordinate creates a new coordinate with the given distance from the
// origin. This should only be used for tests.
func GenerateCoordinate(rtt time.Duration) *coordinate.Coordinate {
	coord := coordinate.NewCoordinate(coordinate.DefaultConfig())
	coord.Vec[0] = rtt.Seconds()
	coord.Height = 0
	return coord
}
