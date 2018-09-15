// ddshop is an implementation of todo-backend [0] that stores data in
// PostgreSQL or CockroachDB.
//
// [0]: https://todobackend.com
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

const listenAddr = ":26256"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "ddshop: %s\n", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	var useCockroach, usePostgres bool
	flagSet := flag.NewFlagSet("ddshop", flag.ExitOnError)
	flagSet.BoolVar(&useCockroach, "cockroach", false,
		"connect to cockroach on ports 26257, 26258, and 26259")
	flagSet.BoolVar(&usePostgres, "postgres", false, "connect to postgres on port 5432")
	flagSet.Usage = func() {
		fmt.Fprintln(os.Stderr, "ddshop: [options] [database-url]...")
		fmt.Fprintln(os.Stderr)
		fmt.Fprintln(os.Stderr, "Options:")
		flagSet.PrintDefaults()
	}
	flagSet.Parse(args)

	if useCockroach && usePostgres {
		return errors.New("-cockroach and -postgres are mutually exclusive")
	} else if useCockroach && flag.NArg() > 0 {
		return errors.New("cannot specify URLs and -cockroach")
	} else if usePostgres && flag.NArg() > 0 {
		return errors.New("cannot specify URLs and -postgres")
	} else if !useCockroach && !usePostgres && flag.NArg() == 0 {
		return errors.New("must specify at least one URL or -cockroach or -postgres")
	}

	urls := flag.Args()
	if useCockroach {
		urls = []string{
			"postgres://root@:26257?sslmode=disable",
			"postgres://root@:26258?sslmode=disable",
			"postgres://root@:26259?sslmode=disable",
		}
	} else if usePostgres {
		urls = []string{"postgres://localhost:5432?sslmode=disable"}
	}

	db, err := connectDB(urls)
	if err != nil {
		return err
	}

	if err := bootstrapDB(db); err != nil {
		return err
	}

	log.Printf("ddshop listening on %s", listenAddr)
	http.Handle("/", &server{db: db})
	return http.ListenAndServe(listenAddr, nil)
}
