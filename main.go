package main

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	cfg = configNew()
)

func main() {
	// Set DEBUG level
	cfg.Log.SetLevel(5)

	// Runing http server
	router := mux.NewRouter().StrictSlash(true)
	router.Use(Middleware)
	//router.HandleFunc("/api/ami", signHandler).Methods("POST")
	//router.HandleFunc("/api/ari", robotHandler).Methods("POST")
	router.HandleFunc("/api/cdr", cdrHandler).Methods("GET")
	http.Handle("/", router)

	srv := http.Server{
		Handler:        router,
		Addr:           cfg.Server.Host + ":" + cfg.Server.Port,
		WriteTimeout:   cfg.Server.Timeout.Write * time.Second,
		ReadTimeout:    cfg.Server.Timeout.Read * time.Second,
		IdleTimeout:    cfg.Server.Timeout.Idle * time.Second,
		MaxHeaderBytes: cfg.Server.MaxHeaderBytes,
		//TLSConfig *tls.Config
	}
	cfg.Log.Println("Server start on " + cfg.Server.Host + ":" + cfg.Server.Port + "...")
	cfg.Log.Fatal(srv.ListenAndServe())
}
