package hansip

import (
	"math/rand"
	"sync"
	"time"
)

type connectionManager struct {
	primary        *connection
	replicas       []*connection
	activeReplicas []*connection
	mut            sync.RWMutex

	connCheckDelay time.Duration

	closed   bool
	quitChan chan struct{}
}

func newConnectionManager(primary *connection, replicas []*connection, connCheckDelay time.Duration) *connectionManager {
	manager := &connectionManager{
		primary:        primary,
		replicas:       replicas,
		mut:            sync.RWMutex{},
		connCheckDelay: connCheckDelay,
		quitChan:       make(chan struct{}),
	}
	manager.updateActiveReplicas()
	go manager.loop()
	return manager
}

func (m *connectionManager) loop() {
	ticker := time.NewTicker(m.connCheckDelay)
	for {
		select {
		case <-ticker.C:
			m.updateActiveReplicas()
		case <-m.quitChan:
			return
		}
	}
}

func (m *connectionManager) getReplicas() []*connection {
	m.mut.RLock()
	slaves := m.replicas
	m.mut.RUnlock()
	return slaves
}

func (m *connectionManager) getActiveReplicas() []*connection {
	m.mut.RLock()
	slaves := m.activeReplicas
	m.mut.RUnlock()
	return slaves
}

func (m *connectionManager) setActiveReplicas(replicas []*connection) {
	m.mut.Lock()
	m.activeReplicas = replicas
	m.mut.Unlock()
}

func (m *connectionManager) updateActiveReplicas() {
	current := m.getReplicas()
	if len(current) == 0 {
		return
	}

	replicas := make([]*connection, 0, len(current))
	for _, conn := range current {
		if conn.isConnected() {
			replicas = append(replicas, conn)
		}
	}

	m.setActiveReplicas(replicas)
}

func (m *connectionManager) getReplica() SQL {
	current := m.getActiveReplicas()
	n := len(current)
	if n == 0 {
		return m.getPrimary()
	}
	return current[rand.Intn(n)].s
}

func (m *connectionManager) getPrimary() SQL {
	if !m.primary.isConnected() {
		return nil
	}
	return m.primary.s
}

func (m *connectionManager) quit() {
	if m.closed {
		return
	}

	m.quitChan <- struct{}{}

	if m.primary != nil {
		m.primary.quit()
	}

	for _, conn := range m.replicas {
		conn.quit()
	}

	m.updateActiveReplicas()
	m.closed = true
}
