package main

import (
	"encoding/json"
	"net/http"
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

func cdrHandler(w http.ResponseWriter, r *http.Request) {

	sd := r.FormValue("startDate")
	ed := r.FormValue("endDate")

	//condition := `WHERE calldate between '` + sd + ` 00:00:00' and '` + ed + ` 23:59:59'`
	condition := `WHERE calldate between ` + sd + ` and ` + ed
	res := GetStatByDate(condition)

	dt, err := json.Marshal(res)
	if err != nil {
		cfg.Log.Errorf("Marshal: %v", err)
		http.Error(w, "marshal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(dt)
}
