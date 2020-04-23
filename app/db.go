package app

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"io/ioutil"
	"strings"
)

type database struct {
	*sql.DB
}

func newDatabase() (database, error) {
	db, err := sql.Open("mysql", dataSourceName)
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
