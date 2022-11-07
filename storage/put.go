package storage

import (
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

func put(w http.ResponseWriter, r *http.Request) {
	root := "/Users/lqy007700/Data/storage"
	create, err := os.Create(root + "/objects/" + strings.Split(r.URL.EscapedPath(), "/")[2])
	defer create.Close()
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	io.Copy(create, r.Body)
}
