package mod

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

//HandlerARI asterisk restful interface handler
func HandlerARI(w http.ResponseWriter, r *http.Request) {
	cfg.Log = cfg.Log.WithFields(logrus.Fields{"mod": "ari", "func": "HandlerARI"})

	w.Write([]byte("ok"))
}
