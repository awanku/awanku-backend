package hansip

import (
	"errors"
	"log"
	"sync/atomic"
	"time"
)

var errPingTimeout = errors.New("ping timeout")

type connection struct {
	s SQL

	pingTimeout    time.Duration
	connCheckDelay time.Duration

	pingFn      func() error
	pingRunning int32
	closeFn     func()

	// 1 for connected, otherwise 0
	connected int32

	closed   bool
	quitChan chan struct{}
}

func newConnection(s SQL, pingTimeout, connCheckDelay time.Duration) *connection {
	conn := &connection{}

	conn.s = s
	conn.pingTimeout = pingTimeout
	conn.connCheckDelay = connCheckDelay
	conn.quitChan = make(chan struct{})
	conn.pingFn = func() error {
		err := s.Exec("select 1;")
		return err
	}
	conn.closeFn = func() {
		err := s.Close()
		if err != nil {
			log.Println("got error when closing database connection:", err)
		}
	}

	// start main loop
	conn.updateStatus()
	go conn.loop()

	return conn
}

func (c *connection) ping() error {
	if !atomic.CompareAndSwapInt32(&c.pingRunning, 0, 1) {
		return nil
	}
	defer atomic.StoreInt32(&c.pingRunning, 0)

	errChan := make(chan error)
	go func() {
		errChan <- c.pingFn()
	}()

	select {
	case <-time.After(c.pingTimeout):
		return errPingTimeout
	case err := <-errChan:
		return err
	}
}

func (c *connection) isConnected() bool {
	return atomic.LoadInt32(&c.connected) == 1
}

func (c *connection) setConnected(connected bool) {
	if connected {
		atomic.StoreInt32(&c.connected, 1)
	} else {
		atomic.StoreInt32(&c.connected, 0)
	}
}

func (c *connection) loop() {
	ticker := time.NewTicker(c.connCheckDelay)
	for {
		select {
		case <-ticker.C:
			c.updateStatus()
		case <-c.quitChan:
			return
		}
	}
}

func (c *connection) updateStatus() {
	connected := c.ping() == nil
	c.setConnected(connected)
}

func (c *connection) quit() {
	if c.closed {
		return
	}

	// stop loop
	c.quitChan <- struct{}{}

	if c.closeFn != nil {
		c.closeFn()
	}

	c.setConnected(false)
	c.closed = true
}
