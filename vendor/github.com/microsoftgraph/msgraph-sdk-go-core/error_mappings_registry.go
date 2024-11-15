package msgraphgocore

import (
	"errors"
	abstractions "github.com/microsoft/kiota-abstractions-go"
	"sync"
)

var lock = &sync.Mutex{}

type errorRegistry struct {
	registry map[string]abstractions.ErrorMappings
}

var singleInstance *errorRegistry

// Create a global thread safe singleton for global values
func getInstance() *errorRegistry {
	if singleInstance == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleInstance == nil {
			singleInstance = &errorRegistry{
				registry: make(map[string]abstractions.ErrorMappings),
			}
		}
	}

	return singleInstance
}

func RegisterError(key string, value abstractions.ErrorMappings) error {
	single := getInstance()
	_, found := single.registry[key]
	if !found {
		single.registry[key] = value
		return nil
	} else {
		return errors.New("object Factory already register")
	}
}

func DeRegisterError(key string) error {
	single := getInstance()
	_, found := single.registry[key]
	if found {
		delete(single.registry, key)
		return nil
	} else {
		return errors.New("object Factory does not exist register")
	}
}

func GetErrorFactoryFromRegistry(key string) (abstractions.ErrorMappings, bool) {
	single := getInstance()
	item, found := single.registry[key]
	return item, found
}
