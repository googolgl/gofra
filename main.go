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
	//cfg.Log.SetLevel(5)

	// init mux router
	router := mux.NewRouter().StrictSlash(true)

	// Runing ARI
	if cfg.ARI.Enable {
		/*if err := ariRun(); err != nil {
			cfg.Log.Panicf("Run ARI: %v", err)
		}*/
		router.HandleFunc("/api/ari/{cmd}", ARIHandler).Methods("POST")
	}

	// Runing AMI
	if cfg.AMI.Enable {
		if _, err := amiRun(); err != nil {
			cfg.Log.Panicf("Run AMI: %v", err)
		}
		router.HandleFunc("/api/ami/{cmd}", AMIHandler).Methods("POST")
	}

	// Runing http server
	router.Use(Middleware)
	router.PathPrefix("/file/").Handler(http.StripPrefix("/file", http.FileServer(http.Dir(cfg.FilePath))))
	router.HandleFunc("/api/cdr", CDRHandler).Methods("GET")
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
