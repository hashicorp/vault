package healthcheck

import (
	"fmt"

	"github.com/hashicorp/vault/sdk/logical"
)

func StringList(source interface{}) ([]string, error) {
	if source == nil {
		return nil, nil
	}

	if value, ok := source.([]string); ok {
		return value, nil
	}

	if rValues, ok := source.([]interface{}); ok {
		var result []string
		for index, rValue := range rValues {
			value, ok := rValue.(string)
			if !ok {
				return nil, fmt.Errorf("unknown source type for []string coercion at index %v: %T", index, rValue)
			}

			result = append(result, value)
		}

		return result, nil
	}

	return nil, fmt.Errorf("unknown source type for []string coercion: %T", source)
}

func fetchMountTune(e *Executor, versionError func()) (*PathFetch, error) {
	tuneRet, err := e.FetchIfNotFetched(logical.ReadOperation, "/sys/mounts/{{mount}}/tune")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch mount tune information: %w", err)
	}

	if !tuneRet.IsSecretOK() {
		if tuneRet.IsUnsupportedPathError() {
			versionError()
		}
	}

	return tuneRet, nil
}
