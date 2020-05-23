package app

import (
	. "ElasticJury/app/common"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"strings"
)

type database struct {
	*sql.DB
}

func newDatabase(databaseName, password string) (database, error) {
	dataSourceName := strings.Replace(strings.Replace(DataSourceName, "<password>", password, 1), "<database>", databaseName, 1)
	fmt.Printf("Using data source: %s\n", dataSourceName)
	db, err := sql.Open("mysql", dataSourceName)
	if err == nil {
		db.SetConnMaxLifetime(DBConnMaxLifeTime)
		err = db.Ping() // test if actually connected
	}
	return database{db}, err
}

// Execute a db command, panic if error occurs
func (db database) mustExec(command string, args ...interface{}) {
	if _, err := db.Exec(command, args...); err != nil {
		panic(err)
	}
}

// Execute db commands from script file
func (db database) execScriptFile(scriptPath string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	bytes, err := ioutil.ReadFile(scriptPath)
	if err != nil {
		return err
	}
	scripts := strings.Split(string(bytes), ";")
	for _, script := range scripts {
		if strings.TrimSpace(script) == "" { // omit null command
			continue
		}
		if _, err := tx.Exec(script); err != nil { // execute init sql command
			return err
		}
	}
	return tx.Commit()
}

// Execute db commands from script file, panic if error occurs
func (db database) mustExecScriptFile(scriptPath string) {
	if err := db.execScriptFile(scriptPath); err != nil {
		panic(err)
	}
}
