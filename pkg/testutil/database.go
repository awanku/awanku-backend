package testutil

import (
	"os"
	"sync"
	"time"

	hansip "github.com/asasmoyo/pq-hansip"
	"github.com/go-pg/pg/v9"
)

var wg sync.WaitGroup

func DBCluster() (*hansip.Cluster, func()) {
	waitChan := make(chan struct{})
	go func() {
		wg.Wait()
		waitChan <- struct{}{}
	}()
	select {
	case <-waitChan:
	case <-time.After(5 * time.Minute):
		panic("failed to get db cluster")
	}

	opts, err := pg.ParseURL(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic("invalid DATABASE_URL")
	}
	db := pg.Connect(opts)
	cluster := hansip.NewCluster(&hansip.Config{
		Primary:        hansip.WrapGoPG(db),
		Replicas:       []hansip.SQL{hansip.WrapGoPG(db)},
		PingTimeout:    1 * time.Second,
		ConnCheckDelay: 5 * time.Second,
	})
	done := func() {
		wg.Done()
	}

	wg.Add(1)
	return cluster, done
}
