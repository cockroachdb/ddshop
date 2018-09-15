package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/cockroachdb/ddshop/robustdb"
	"github.com/elazarl/go-bindata-assetfs"
)

//go:generate go-bindata -prefix ui/build ui/build

var cwd = func() string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return cwd
}()

type server struct {
	db         *robustdb.DB
	fileServer http.Handler
}

func newServer(db *robustdb.DB, dev bool) *server {
	s := &server{
		db: db,
	}
	if dev {
		log.Printf("using assets on disk")
		s.fileServer = http.FileServer(http.Dir("ui/build"))
	} else {
		log.Printf("using assets embedded in the binary")
		s.fileServer = http.FileServer(&assetfs.AssetFS{
			Asset:     Asset,
			AssetDir:  AssetDir,
			AssetInfo: AssetInfo,
		})
	}
	return s
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
	s.fileServer.ServeHTTP(w, r)
}

func parseTodoID(r *http.Request) (int32, error) {
	parts := strings.SplitN(r.URL.Path, "/", 3)
	if n := len(parts); n != 3 {
		return 0, fmt.Errorf("expected 3 path segments, got %d", n)
	}
	id, err := strconv.Atoi(parts[2])
	return int32(id), err
}

func writeJSON(w http.ResponseWriter, data interface{}) error {
	bytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	return err
}

func (s *server) serveAPI(w http.ResponseWriter, r *http.Request, body []byte) {
	switch r.Method {
	case "GET":
		todos, err := listTodos(s.db)
		if err != nil {
			writeError(w, err)
			return
		}
		if err := writeJSON(w, todos); err != nil {
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
		if err := writeJSON(w, &t); err != nil {
			writeError(w, err)
			return
		}
	case "DELETE":
		id, err := parseTodoID(r)
		if err != nil {
			writeError(w, err)
			return
		}
		if err := deleteTodo(s.db, id); err != nil {
			writeError(w, err)
			return
		}
		if err := writeJSON(w, &todo{}); err != nil {
			writeError(w, err)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("forbidden HTTP method %s", r.Method), http.StatusMethodNotAllowed)
	}
}
