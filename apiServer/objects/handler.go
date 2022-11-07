package objects

import "net/http"

func Handler(w http.ResponseWriter, r *http.Request)  {
	m := r.Method
	switch m {
	case http.MethodGet:

	case http.MethodPut:
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)

	}



}
