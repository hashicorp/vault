package agent

import (
	"sync"

	"github.com/hashicorp/vault/api"
)

func NewClientManager(rootClient *api.Client) *ClientManager {
	return &ClientManager{
		rootClient: rootClient,
		clients:    make(map[*api.Client]struct{}),
		l:          &sync.RWMutex{},
	}
}

type ClientManager struct {
	rootClient *api.Client
	clients    map[*api.Client]struct{}
	l          *sync.RWMutex
}

func (m *ClientManager) New() (*api.Client, error) {
	m.l.Lock()
	defer m.l.Unlock()

	client, err := m.rootClient.Clone()
	if err != nil {
		return nil, err
	}

	m.clients[client] = struct{}{}
	return client, nil
}

func (m *ClientManager) SetToken(token string) {
	m.l.Lock()
	defer m.l.Unlock()

	m.rootClient.SetToken(token)

	for c, _ := range m.clients {
		c.SetToken(token)
	}
}

func (m *ClientManager) Remove(client *api.Client) {
	m.l.Lock()
	defer m.l.Unlock()
	delete(m.clients, client)
}
