package db

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/linuxerwang/confish"
	db "upper.io/db.v3"
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/mysql"
	"upper.io/db.v3/postgresql"
)

var (
	dbConfFile = flag.String("db-conf", "/opt/skylb-dashboard/conf/database.conf", "The database config file")
	dbClient   sqlbuilder.Database

	debugMode bool
)

type databaseConf struct {
	DBType       string `cfg-attr:"dbtype"`
	Host         string `cfg-attr:"host"`
	Port         int32  `cfg-attr:"port"`
	User         string `cfg-attr:"user"`
	Password     string `cfg-attr:"password"`
	DbName       string `cfg-attr:"db-name"`
	MaxIdleConns int    `cfg-attr:"max-idle-connections"`
	MaxOpenConns int    `cfg-attr:"max-open-connections"`
}

type confWrapper struct {
	DbConf *databaseConf `cfg-attr:"database"`
}

func checkFlags() {
	if *dbConfFile == "" {
		fmt.Println("Flag db-conf is required.")
		os.Exit(2)
	}

	if !debugMode && *ldapEndpoint == "" {
		fmt.Println("Flag --ldap-endpoint is required.")
	}
}

// Init initializes the database.
func Init(debug bool) {
	debugMode = debug

	checkFlags()

	f, err := os.Open(*dbConfFile)
	if err != nil {
		fmt.Println("Failed to open the database config file,", err)
	}
	defer f.Close()

	dbClient = initFromReader(f)
}

// Example database config:
//
// database {
//     host: localhost
//     port: 5432
//     user: skylb_dashboard
//     password: passwd
//     db-name: skylb_dashboard
// }
func initFromReader(r io.Reader) sqlbuilder.Database {
	confWrapper := &confWrapper{}
	confish.Parse(r, confWrapper)
	dbConf := confWrapper.DbConf

	// Validate config values.
	if dbConf.Host == "" || dbConf.User == "" || dbConf.Password == "" || dbConf.DbName == "" {
		log.Fatalln("Not enough database config information.")
	}
	if dbConf.Port == 0 {
		dbConf.Port = 5432
	}
	if dbConf.MaxIdleConns == 0 {
		dbConf.MaxIdleConns = 5
	}
	if dbConf.MaxOpenConns == 0 {
		dbConf.MaxOpenConns = 50
	}

	db.DefaultSettings.SetMaxIdleConns(dbConf.MaxIdleConns)
	db.DefaultSettings.SetMaxOpenConns(dbConf.MaxOpenConns)

	var settings db.ConnectionURL
	switch dbConf.DBType {
	case "", "postgresql":
		settings = postgresql.ConnectionURL{
			Database: dbConf.DbName,
			Host:     dbConf.Host,
			User:     dbConf.User,
			Password: dbConf.Password,
		}
	case "mysql":
		settings = mysql.ConnectionURL{
			Database: dbConf.DbName,
			Host:     dbConf.Host,
			User:     dbConf.User,
			Password: dbConf.Password,
		}
	}

	var session sqlbuilder.Database
	var err error
	for {
		switch dbConf.DBType {
		case "", "postgresql":
			session, err = postgresql.Open(settings)
		case "mysql":
			session, err = mysql.Open(settings)
		}
		if err == nil {
			return session
		}

		log.Printf("Failed to open database, %+v. Wait for one second.\n", err)
		time.Sleep(time.Second)
	}
}
