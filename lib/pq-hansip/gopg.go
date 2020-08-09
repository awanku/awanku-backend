package hansip

import (
	pg "github.com/go-pg/pg/v9"
)

type gopgSQL struct {
	db *pg.DB
}

func WrapGoPG(db *pg.DB) SQL {
	return &gopgSQL{
		db: db,
	}
}

func (s *gopgSQL) Query(dest interface{}, query string, args ...interface{}) error {
	query = injectCallerInfo(query)
	_, err := s.db.Query(dest, query, args...)
	return err
}

func (s *gopgSQL) Exec(query string, args ...interface{}) error {
	query = injectCallerInfo(query)
	_, err := s.db.Exec(query, args...)
	return err
}

func (s *gopgSQL) NewTransaction() (Transaction, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	return &gopgTransaction{db: tx}, nil
}

func (s *gopgSQL) Close() error {
	return s.db.Close()
}

type gopgTransaction struct {
	db       *pg.Tx
	finished bool
}

func (tx *gopgTransaction) Query(dest interface{}, query string, args ...interface{}) error {
	query = injectCallerInfo(query)
	_, err := tx.db.Query(dest, query, args...)
	return err
}

func (tx *gopgTransaction) Exec(query string, args ...interface{}) error {
	query = injectCallerInfo(query)
	_, err := tx.db.Exec(query, args...)
	return err
}

func (tx *gopgTransaction) Commit() error {
	if tx.finished {
		return ErrTxFinished
	}
	err := tx.db.Commit()
	tx.finished = true
	return err
}

func (tx *gopgTransaction) Rollback() error {
	if tx.finished {
		return ErrTxFinished
	}
	err := tx.db.Rollback()
	tx.finished = true
	return err
}
