package handlers

import (
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"net/http"
	"path/filepath"
)

// Files is a handler for reading and writing files
type Files struct {
	log   hclog.Logger
	store files.Storage
}

// NewFiles create a new File handler
func NewFiles(s files.Storage, l hclog.Logger) *Files {
	return &Files{store: s, log: l}
}

// ServeHTTP implements the http.Handler interface
func (f *Files) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	f.log.Info("Handle POST", "id", id, "filename", fn)

	// check that the filepath is a valid name and file
	if id == "" || fn == "" {
		f.invalidURI(r.URL.String(), rw)
		return
	}

	f.saveFile(id, fn, rw, r)
}

func (f *Files) invalidURI(uri string, rw http.ResponseWriter) {
	f.log.Error("Invalid", "path", uri)
	http.Error(rw, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}

// saveFile saves the contents of the request to a file
func (f *Files) saveFile(id, path string, rw http.ResponseWriter, r *http.Request) {
	f.log.Info("Save file for product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r.Body)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}
