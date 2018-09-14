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
	log.Printf("%s %s", r.Method, r.URL)
	if !strings.HasPrefix(r.URL.Path, "/api") {
		http.ServeFile(w, r, filepath.Join(cwd, "assets", r.URL.Path))
		return
	} else if r.URL.Path != "/api" {
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

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
	case "POST":
		var t todo
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			writeError(w, err)
			return
		}
		if err := json.Unmarshal(body, &t); err != nil {
			writeError(w, err)
			return
		}
		if err := upsertTodo(s.db, t); err != nil {
			writeError(w, err)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("forbidden HTTP method %s", r.Method), http.StatusMethodNotAllowed)
	}
}
