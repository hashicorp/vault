package auth

import (
	"crypto/tls"
	"errors"
	"github.com/oracle/oci-go-sdk/common"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

var customTransport = &http.Transport{
	TLSClientConfig: &tls.Config{
		ServerName: "test",
	},
}

func TestNewDispatcherModifier_NoInitialModifier(t *testing.T) {
	modifier := newDispatcherModifier(nil)
	initialClient := &http.Client{}
	returnClient, err := modifier.Modify(initialClient)
	assert.Nil(t, err)
	assert.ObjectsAreEqual(initialClient, returnClient)
}

func TestNewDispatcherModifier_InitialModifier(t *testing.T) {
	modifier := newDispatcherModifier(setCustomCAPool)
	initialClient := &http.Client{}
	returnDispatcher, err := modifier.Modify(initialClient)
	assert.Nil(t, err)
	returnClient := returnDispatcher.(*http.Client)
	assert.ObjectsAreEqual(returnClient.Transport, customTransport)
}

func TestNewDispatcherModifier_ModifierFails(t *testing.T) {
	modifier := newDispatcherModifier(modifierGoneWrong)
	initialClient := &http.Client{}
	returnClient, err := modifier.Modify(initialClient)
	assert.NotNil(t, err)
	assert.Nil(t, returnClient)
}

func setCustomCAPool(dispatcher common.HTTPRequestDispatcher) (common.HTTPRequestDispatcher, error) {
	client := dispatcher.(*http.Client)
	client.Transport = customTransport
	return client, nil
}

func modifierGoneWrong(dispatcher common.HTTPRequestDispatcher) (common.HTTPRequestDispatcher, error) {
	return nil, errors.New("uh oh")
}
