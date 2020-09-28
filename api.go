package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

//Middleware - main handler
func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > 1024 {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			cfg.Log.Error("request body too large")
			return
		}

		next.ServeHTTP(w, r)
	})
}

//ARIHandler asterisk restful interface handler
func ARIHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("ok"))
}

//AMIHandler asterisk manager interface handler
func AMIHandler(w http.ResponseWriter, r *http.Request) {
	switch mux.Vars(r)["cmd"] {
	case "my":

	}

	w.Write([]byte("ok"))
}

//CDRHandler asterisk call detail records handler
func CDRHandler(w http.ResponseWriter, r *http.Request) {
	var respData []byte

	sd := r.FormValue("startDate")
	ed := r.FormValue("endDate")

	if len(sd) > 0 && len(ed) > 0 {
		condition := `WHERE calldate between ` + sd + ` and ` + ed
		res, err := GetStatBy(condition)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		respData, err = json.Marshal(res)
		if err != nil {
			cfg.Log.Errorf("marshal: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if respData == nil {
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
	}

	w.Write(respData)
}
