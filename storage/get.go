package storage

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func get(w http.ResponseWriter, r *http.Request) {
	root := "/Users/lqy007700/Data/storage"
	open, err := os.Open(root + "/objects/" + strings.Split(r.URL.EscapedPath(), "/")[2])
	defer open.Close()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	io.Copy(w,open)
}
