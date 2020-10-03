package mod

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

// CDR structure
type CDR struct {
	CallDate      time.Time `json:"calldate"`
	Src           string    `json:"src"`
	Dst           string    `json:"dst"`
	Dcontext      string    `json:"dcontext"`
	Channel       string    `json:"channel"`
	Disposition   string    `json:"disposition"`
	Dstchannel    string    `json:"dstchannel"`
	Lastapp       string    `json:"lastapp"`
	Duration      int       `json:"duration"`
	Billsec       int       `json:"billsec"`
	Sequence      int       `json:"sequence"`
	UniqueID      string    `json:"uniqueid"`
	ActionID      string    `json:"actionid"`
	RecordingFile string    `json:"recordingfile,omitempty"`
}

//HandlerCDR asterisk call detail records handler
func HandlerCDR(w http.ResponseWriter, r *http.Request) {
	cfg.Log = cfg.Log.WithFields(logrus.Fields{"mod": "cdr", "func": "HandlerCDR"})

	var respData []byte

	sd := r.FormValue("StartDate")
	ed := r.FormValue("EndDate")

	if len(sd) > 0 && len(ed) > 0 {
		condition := fmt.Sprintf(`WHERE calldate between %s and %s`, sd, ed)
		res, err := cdrGetStatBy(condition)
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
		cfg.Log.Debug("Invalid keys")
		w.WriteHeader(http.StatusBadRequest)
	} else {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
	}

	w.Write(respData)
}

//cdrGetStatBy - main database function
func cdrGetStatBy(condition string) ([]CDR, error) {
	db, err := connectDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	cfg.Log = cfg.Log.WithFields(logrus.Fields{"mod": "cdr", "func": "cdrGetStatBy"})

	queryString := `SELECT calldate, src, dst, dcontext, channel, dstchannel, ` +
		`lastapp, duration, billsec, disposition, uniqueid, sequence FROM cdr ` + condition

	if cfg.CDR.Recname != "" {
		queryString = `SELECT calldate, src, dst, dcontext, channel, dstchannel, ` +
			`lastapp, duration, billsec, disposition, uniqueid, recordingfile, sequence ` +
			`FROM cdr ` + condition
	}

	rows, err := db.Query(queryString)
	if err != nil {
		cfg.Log.Errorf("query: %v", err)
		return nil, err
	}
	defer rows.Close()

	var cdr CDR
	var arrayCDR []CDR
	// Fetch rows
	for rows.Next() {
		if cfg.CDR.Recname != "" {
			err = rows.Scan(&cdr.CallDate, &cdr.Src, &cdr.Dst, &cdr.Dcontext, &cdr.Channel,
				&cdr.Dstchannel, &cdr.Lastapp, &cdr.Duration, &cdr.Billsec, &cdr.Disposition,
				&cdr.UniqueID, &cdr.RecordingFile, &cdr.Sequence)
		} else {
			err = rows.Scan(&cdr.CallDate, &cdr.Src, &cdr.Dst, &cdr.Dcontext, &cdr.Channel,
				&cdr.Dstchannel, &cdr.Lastapp, &cdr.Duration, &cdr.Billsec, &cdr.Disposition,
				&cdr.UniqueID, &cdr.Sequence)
		}
		if err != nil {
			cfg.Log.Errorf("fetch: %v", err)
			break
		}
		arrayCDR = append(arrayCDR, cdr)
	}

	if err = rows.Err(); err != nil {
		cfg.Log.Errorf("error handling: %v", err)
		return nil, err
	}

	return arrayCDR, nil
}
