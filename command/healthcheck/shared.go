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

func fetchMountTune(e *Executor, versionError func()) (bool, *PathFetch, map[string]interface{}, error) {
	tuneRet, err := e.FetchIfNotFetched(logical.ReadOperation, "/sys/mounts/{{mount}}/tune")
	if err != nil {
		return true, nil, nil, err
	}

	if !tuneRet.IsSecretOK() {
		if tuneRet.IsUnsupportedPathError() {
			versionError()
		}

		return true, nil, nil, nil
	}

	var data map[string]interface{} = nil
	if len(tuneRet.Secret.Data) > 0 {
		data = tuneRet.Secret.Data
	}

	return false, tuneRet, data, nil
}
