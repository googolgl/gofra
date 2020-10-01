package mod

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

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
	defer clientAMI.Close()

	var respData *gami.AMIResponse
	switch typeAction {
	case "sync":
		respData, err = clientAMI.Action(pami)
		if err != nil {
			cfg.Log.Errorf("action ami: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	case "async":
		//pami["ActionID"] = fmt.Sprintf("%d", time.Now().UnixNano())
		_, err := clientAMI.AsyncAction(pami)
		if err != nil {
			cfg.Log.Errorf("asyncAction ami: %v", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		respData = &gami.AMIResponse{
			ID:     pami["ActionID"],
			Status: "Success",
		}
	}

	rsp, err := json.Marshal(respData)
	if err != nil {
		cfg.Log.Errorf("marshal : %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(rsp)
}
