package mod

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/googolgl/gami"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

//ActionAMI - ami action structure
type ActionAMI struct {
	Action      string `json:"Action"`
	Channel     string `json:"Channel"`
	Application string `json:"Application,omitemti"`
	Context     string `json:"Context,omitemti"`
	Exten       string `json:"Exten,omitemti"`
	Data        string `json:"Data,omitemti"`
	Priority    int    `json:"Priority,omitemti"`
	Callerid    int    `json:"Callerid,omitemti"`
	Variable    string `json:"Variable,omitemti"`
}

func connectAMI() (*gami.AMIClient, error) {
	amiClient, err := gami.Dial(cfg.AMI.Host + ":" + cfg.AMI.Port)
	if err != nil {
		return nil, err
	}
	//defer amiClient.Close()

	amiClient.Run()

	if err := amiClient.Login(cfg.AMI.UserName, cfg.AMI.Password); err != nil {
		return nil, err
	}

	return amiClient, nil
}

//HandlerAMI asterisk manager interface handler
func HandlerAMI(w http.ResponseWriter, r *http.Request) {
	cfg.Log = cfg.Log.WithFields(logrus.Fields{"mod": "ami", "func": "HandlerAMI"})

	typeAction := mux.Vars(r)["type"]
	if typeAction != "sync" && typeAction != "async" {
		cfg.Log.Errorf("not found: %v", typeAction)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		cfg.Log.Errorf("read body: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var pami gami.Params
	if err = json.Unmarshal(body, &pami); err != nil {
		cfg.Log.Errorf("unmarshal data: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	clientAMI, err := connectAMI()
	if err != nil {
		cfg.Log.Errorf("connection to ami: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var rsp []byte
	var respData *gami.AMIResponse
	switch typeAction {
	case "sync":
		respAction, _, err := clientAMI.Action(pami)
		if err != nil {
			cfg.Log.Errorf("action ami: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer clientAMI.Close()

		respData = (<-respAction)

	case "async":
		webHook, _ := pami["webhook"]

		respAction, respActionID, err := clientAMI.Action(pami)
		if err != nil {
			cfg.Log.Errorf("action ami: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		go func() {
			defer clientAMI.Close()

			if len(webHook) != 0 {
				// Send webhook
				rsp, err = json.Marshal(<-respAction)
				if err != nil {
					cfg.Log.Errorf("marshal : %v", err)
					return
				}

				url, err := url.Parse(webHook)
				if err != nil {
					cfg.Log.Errorf("url parsing: %v", err)
					return
				}

				resp, err := http.Post(url.String(), "application/json", bytes.NewBuffer(rsp))
				if err != nil {
					cfg.Log.Errorf("post request: %v", err)
					return
				}
				defer resp.Body.Close()
			}

		}()

		respData = &gami.AMIResponse{
			ID:     respActionID,
			Status: "Accepted",
		}
	}

	rsp, err = json.Marshal(respData)
	if err != nil {
		cfg.Log.Errorf("marshal : %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(rsp)
}
