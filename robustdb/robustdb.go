// Package robustdb implements a database client that can tolerate failures.
package robustdb

import (
	"context"
	"database/sql"
)

type DB struct {
	urls  []string
	conns []*sql.Conn
}

func New(urls ...string) (*DB, error) {
	ctx := context.TODO()
	db := &DB{urls: urls}
	for _, u := range urls {
		sqlDB, err := sql.Open("postgres", u)
		if err != nil {
			return nil, err
		}
		conn, err := sqlDB.Conn(ctx)
		if err != nil {
			return nil, err
		}
		if _, err := conn.ExecContext(ctx, `SET statement_timeout = '200ms'`); err != nil {
			return nil, err
		}
		db.conns = append(db.conns, conn)
	}
	return db, nil
}

func (db *DB) retry(fn func(conn *sql.Conn) error) error {
	var err error
	for _, sqlDB := range db.conns {
		if err = fn(sqlDB); err == nil {
			return err
		}
	}
	return err
}

func (db *DB) Ping() error {
	return db.retry(func(conn *sql.Conn) error {
		return conn.PingContext(context.TODO())
	})
}

func (db *DB) Exec(query string, args ...interface{}) (sql.Result, error) {
	var result sql.Result
	err := db.retry(func(conn *sql.Conn) error {
		var err error
		result, err = conn.ExecContext(context.TODO(), query, args...)
		return err
	})
	return result, err
}

func (db *DB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	var rows *sql.Rows
	err := db.retry(func(conn *sql.Conn) error {
		var err error
		rows, err = conn.QueryContext(context.TODO(), query, args...)
		return err
	})
	return rows, err
}
