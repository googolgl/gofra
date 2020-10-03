package main

import (
	"net/http"
	"time"

	"github.com/googolgl/gofra/mod"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	cfg = mod.ConfigNew()
)

func main() {
	cfg.Log = cfg.Log.WithFields(logrus.Fields{"mod": "main", "func": "main"})

	// init mux router
	router := mux.NewRouter().StrictSlash(true)

	// Activate ARI
	if cfg.ARI.Enable {
		cfg.Log.Debug("cfg.ARI.Enable")
		router.HandleFunc("/api/ari/{cmd}", mod.HandlerARI).Methods("POST")
	}

	// Activate AMI
	if cfg.AMI.Enable {
		router.HandleFunc("/api/ami/{type}", mod.HandlerAMI).Methods("POST")
	}

	// Activate CDR
	if cfg.CDR.Enable {
		router.HandleFunc("/api/cdr", mod.HandlerCDR).Methods("GET")
	}

	// Activate CEL
	if cfg.CEL.Enable {
		router.HandleFunc("/api/cel", mod.HandlerCDR).Methods("GET")
	}

	// Runing http server
	router.Use(Middleware)
	router.PathPrefix("/file/").Handler(http.StripPrefix("/file", http.FileServer(http.Dir(cfg.FilePath))))
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

	cfg.Log.Warnf("Starting on " + cfg.Server.Host + ":" + cfg.Server.Port + "...")
	cfg.Log.Fatal(srv.ListenAndServe())
}

//Middleware - main handler
func Middleware(next http.Handler) http.Handler {
	cfg.Log = cfg.Log.WithFields(logrus.Fields{"mod": "main", "func": "Middleware"})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > 1024 {
			http.Error(w, "request body too large", http.StatusRequestEntityTooLarge)
			cfg.Log.Warn("request body too large")
			return
		}

		// Autorisation

		next.ServeHTTP(w, r)
	})
}
