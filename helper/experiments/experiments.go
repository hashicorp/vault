package experiments

const VaultExperimentEventsBeta1 = "events.beta1"

var validExperiments = []string{
	VaultExperimentEventsBeta1,
}

// ValidExperiments exposes the list without exposing a mutable global variable.
// Experiments can only be enabled when starting a server, and will typically
// enable pre-GA API functionality.
func ValidExperiments() []string {
	result := make([]string, len(validExperiments))
	copy(result, validExperiments)
	return result
}
