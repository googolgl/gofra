package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
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
	Amaflags      int       `json:"amaflags"`
	Sequence      int       `json:"sequence"`
	UniqueID      string    `json:"uniqueid"`
	ActionID      string    `json:"actionid"`
	RecordingFile string    `json:"recordingfile,omitempty"`
}

func connectDB() (*sql.DB, error) {
	var connString string
	switch cfg.DB.DrvName {
	case "mysql":
		connString = fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			cfg.DB.UserName, cfg.DB.Password, cfg.DB.Host, cfg.DB.Database)

	default:
		cfg.Log.Fatalf("Invalid database driver: %s", cfg.DB.DrvName)
	}

	db, err := sql.Open(cfg.DB.DrvName, connString)
	if err != nil {
		return nil, err
	}

	return db, nil
}

//GetStatByDate - main database function
func GetStatByDate(condition string) []CDR {
	db, err := connectDB()
	if err != nil {
		cfg.Log.Errorf("connect: %v", err)
	}
	defer db.Close()

	rows, err := db.Query(`SELECT calldate, src, dst, dcontext, channel, dstchannel, ` +
		`lastapp, duration, billsec, disposition, amaflags, uniqueid, recordingfile, sequence ` +
		`FROM cdr ` + condition)
	if err != nil {
		cfg.Log.Errorf("query: %v", err)
	}
	defer rows.Close()

	var cdr CDR
	var arrayCDR []CDR
	// Fetch rows
	for rows.Next() {
		err = rows.Scan(&cdr.CallDate, &cdr.Src, &cdr.Dst, &cdr.Dcontext, &cdr.Channel,
			&cdr.Dstchannel, &cdr.Lastapp, &cdr.Duration, &cdr.Billsec, &cdr.Disposition,
			&cdr.Amaflags, &cdr.UniqueID, &cdr.RecordingFile, &cdr.Sequence)
		if err != nil {
			cfg.Log.Errorf("fetch: %v", err)
			break
		}
		arrayCDR = append(arrayCDR, cdr)
	}

	if err = rows.Err(); err != nil {
		cfg.Log.Errorf("error handling: %v", err)
	}
	/*if startdate == "" || enddate == "" {
		query.Table("cdr").Where("src=? or dst=?", MSISDN, MSISDN).Find(&cdrs)
	} else {
		query.Table("cdr").Where("calldate between ? and ? and (src=? or dst=?)",
			startdate+" 00:00:00", enddate+" 23:59:00", MSISDN, MSISDN).Find(&cdrs)
	}*/

	return arrayCDR
}
