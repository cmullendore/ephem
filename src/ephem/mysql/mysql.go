// .sql files are embedded using go-bindata and the following command:
// go-bindata.exe -prefix "src/data/mysql/" -o src/data/mysql/procedures.go src/data/mysql/sql/...

package mysql

import (
	"database/sql"
	"log"
	"time"

	"github.com/cmullendore/ephem/src/config"

	"github.com/go-sql-driver/mysql"
)

type MySQL struct {
	db *sql.DB
	l  *log.Logger
}

func Initialize(c *config.Ephem) *MySQL {
	log.Println("Initializing database connection - MySQL")
	var mysqlConfig = mysql.NewConfig()
	mysqlConfig.MultiStatements = true
	mysqlConfig.ParseTime = true
	mysqlConfig.Addr = c.Database.Server
	mysqlConfig.Net = "tcp"
	mysqlConfig.DBName = c.Database.Name
	mysqlConfig.User = c.Database.Username
	mysqlConfig.Passwd = c.Database.Password

	var mysqlConn, err = mysql.NewConnector(mysqlConfig)
	if err != nil {
		log.Fatalln("Failed to create MySQL connection")
	}

	var db = sql.OpenDB(mysqlConn)

	log.Println("Creating secrets table")

	if _, err = db.Exec(createSecretsTable); err != nil {
		log.Println("Failed to initialize secrets table")
		log.Panicln(err)
	}

	log.Println("Created secrets table")

	sql := MySQL{db: db}

	return &sql
}

func (db *MySQL) SaveItem(path, item *string, lifetime, readCount int) *error {

	var t = time.Now()

	if _, err := db.db.Exec(saveSecret, path, item, t, lifetime, readCount, t); err != nil {
		return &err
	}

	return nil
}

func (db *MySQL) GetItem(path *string) (*string, *error) {
	row := db.db.QueryRow(getSecret, path)

	var secret string

	if err := row.Scan(&secret); err != nil {
		log.Println(err)
		return nil, &err
	}

	return &secret, nil
}

func (db *MySQL) IncrementReadCount(path *string) *error {
	_, err := db.db.Exec(incrementReadCount, path)
	if err != nil {
		log.Println(err)
		return &err
	}

	return nil
}

func (db *MySQL) DeleteItem(path *string) *error {
	_, err := db.db.Exec(deleteSecret, path)
	if err != nil {
		log.Println(err)
		return &err
	}

	return nil
}

func (db *MySQL) CleanupSecrets(maxReads int) *error {
	_, err := db.db.Exec(cleanupSecrets, maxReads)
	if err != nil {
		log.Println(err)
		return &err
	}

	return nil
}
