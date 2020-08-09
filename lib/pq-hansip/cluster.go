package hansip

import (
	"errors"
	"time"
)

var (
	ErrNoReplicaAvailable = errors.New("no replica connection available")
	ErrNoPrimaryAvailable = errors.New("no primary connection available")
)

type Config struct {
	Primary  SQL
	Replicas []SQL

	PingTimeout    time.Duration
	ConnCheckDelay time.Duration
}

type ClusterHealth struct {
	PrimaryOK  bool `json:"primary_ok"`
	ReplicasOK bool `json:"replicas_ok"`
}

type Cluster struct {
	manager *connectionManager
}

func NewCluster(conf *Config) *Cluster {
	primary := newConnection(conf.Primary, conf.PingTimeout, conf.ConnCheckDelay)
	replicas := make([]*connection, 0)
	for _, replica := range conf.Replicas {
		replicas = append(replicas, newConnection(replica, conf.PingTimeout, conf.ConnCheckDelay))
	}

	cluster := Cluster{
		manager: newConnectionManager(
			primary,
			replicas,
			conf.ConnCheckDelay/2,
		),
	}
	return &cluster
}

func (c *Cluster) Query(dest interface{}, query string, args ...interface{}) error {
	conn := c.manager.getReplica()
	if conn == nil {
		return ErrNoReplicaAvailable
	}
	return conn.Query(dest, query, args...)
}

func (c *Cluster) WriterExec(query string, args ...interface{}) error {
	conn := c.manager.getPrimary()
	if conn == nil {
		return ErrNoPrimaryAvailable
	}
	return conn.Exec(query, args...)
}

func (c *Cluster) WriterQuery(dest interface{}, query string, args ...interface{}) error {
	conn := c.manager.getPrimary()
	if conn == nil {
		return ErrNoPrimaryAvailable
	}
	return conn.Query(dest, query, args...)
}

func (c *Cluster) NewTransaction() (Transaction, error) {
	conn := c.manager.getPrimary()
	if conn == nil {
		return nil, ErrNoPrimaryAvailable
	}
	return conn.NewTransaction()
}

func (c *Cluster) Shutdown() {
	c.manager.quit()
}

func (c *Cluster) Health() *ClusterHealth {
	var replicasOK bool
	for _, replica := range c.manager.getActiveReplicas() {
		if replica.isConnected() {
			replicasOK = true
			break
		}
	}
	return &ClusterHealth{
		PrimaryOK:  c.manager.primary.isConnected(),
		ReplicasOK: replicasOK,
	}
}
