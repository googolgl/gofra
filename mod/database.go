package mod

import (
	"database/sql"
	"fmt"

	//use mysql module
	_ "github.com/go-sql-driver/mysql"
	"github.com/sirupsen/logrus"
)

func connectDB() (*sql.DB, error) {
	cfg.Log = cfg.Log.WithFields(logrus.Fields{"mod": "database", "func": "connectDB"})

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
		cfg.Log.Errorf("open database: %v", err)
		return nil, err
	}

	return db, nil
}
