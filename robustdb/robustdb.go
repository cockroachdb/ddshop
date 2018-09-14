// Package robustdb implements a database client that is robust to failures.
package robustdb

import "database/sql"

type DB struct {
	urls   []string
	sqlDBs []*sql.DB
}

func New(urls ...string) (*DB, error) {
	db := &DB{urls: urls}
	for _, u := range urls {
		sqlDB, err := sql.Open("postgres", u)
		if err != nil {
			return nil, err
		}
		db.sqlDBs = append(db.sqlDBs, sqlDB)
	}
	return db, nil
}

func (db *DB) retry(fn func(sqlDB *sql.DB) error) error {
	var err error
	for _, sqlDB := range db.sqlDBs {
		if err = fn(sqlDB); err == nil {
			return err
		}
	}
	return err
}

func (db *DB) Ping() error {
	return db.retry(func(sqlDB *sql.DB) error {
		return sqlDB.Ping()
	})
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	err := db.retry(func(sqlDB *sql.DB) error {
		var err error
		result, err = sqlDB.Exec(query, args...)
		return err
	})
	return result, err
}

func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	var rows *sql.Rows
	err := db.retry(func(sqlDB *sql.DB) error {
		var err error
		rows, err = sqlDB.Query(query, args...)
		return err
	})
	return rows, err
}
