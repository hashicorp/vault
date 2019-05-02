package license

// Features is a bitmask of feature flags
type Features uint

const FeatureNone Features = 0

func (f Features) HasFeature(flag Features) bool {
	return false
}
