package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/cockroachdb/ddshop/robustdb"
)

var cwd = func() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}()

type server struct {
	db *robustdb.DB
}

func writeError(w http.ResponseWriter, err error) {
	http.Error(w, err.Error(), http.StatusInternalServerError)
	log.Printf("error while handling request: %s", err)
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		writeError(w, err)
		return
	}

	// Log request.
	logDetail := ""
	if len(body) > 0 {
		logDetail += " " + string(body)
	}
	log.Printf("%s %s%s", r.Method, r.URL, logDetail)

	if strings.HasPrefix(r.URL.Path, "/api") {
		s.serveAPI(w, r, body)
	} else {
		s.serveFile(w, r)
	}
}

func (s *server) serveFile(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, filepath.Join(cwd, "assets", r.URL.Path))
}

func (s *server) serveAPI(w http.ResponseWriter, r *http.Request, body []byte) {
	switch r.Method {
	case "GET":
		todos, err := listTodos(s.db)
		if err != nil {
			writeError(w, err)
			return
		}
		data, err := json.Marshal(todos)
		if err != nil {
			writeError(w, err)
			return
		}
		if _, err := w.Write(data); err != nil {
			writeError(w, err)
			return
		}
	case "PUT", "POST":
		var t todo
		if err := json.Unmarshal(body, &t); err != nil {
			writeError(w, err)
			return
		}
		if err := upsertTodo(s.db, &t); err != nil {
			writeError(w, err)
			return
		}
		data, err := json.Marshal(&t)
		if err != nil {
			writeError(w, err)
			return
		}
		if _, err := w.Write(data); err != nil {
			writeError(w, err)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("forbidden HTTP method %s", r.Method), http.StatusMethodNotAllowed)
	}
}
