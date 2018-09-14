package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data, err := json.Marshal(todos)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if _, err := w.Write(data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	case "POST":
		var t todo
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(body, &t); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := upsertTodo(s.db, t); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("forbidden HTTP method %s", r.Method), http.StatusMethodNotAllowed)
	}
}
