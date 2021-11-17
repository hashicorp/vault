package jwt

func verifyPrincipals(pcpls, auds []string) bool {
	// "Each principal intended to process the JWT MUST
	// identify itself with a value in the audience claim."
	// - https://tools.ietf.org/html/rfc7519#section-4.1.3

	found := -1
	for i, p := range pcpls {
		for _, v := range auds {
			if p == v {
				found++
				break
			}
		}
		if found != i {
			return false
		}
	}
	return true
}

// ValidAudience returns true iff:
// 	- a and b are strings and a == b
// 	- a is string, b is []string and a is in b
// 	- a is []string, b is []string and all of a is in b
// 	- a is []string, b is string and len(a) == 1 and a[0] == b
func ValidAudience(a, b interface{}) bool {
	s1, ok := a.(string)
	if ok {
		if s2, ok := b.(string); ok {
			return s1 == s2
		}
		a2, ok := b.([]string)
		return ok && verifyPrincipals([]string{s1}, a2)
	}

	a1, ok := a.([]string)
	if !ok {
		return false
	}
	if a2, ok := b.([]string); ok {
		return verifyPrincipals(a1, a2)
	}
	s2, ok := b.(string)
	return ok && len(a1) == 1 && a1[0] == s2
}
