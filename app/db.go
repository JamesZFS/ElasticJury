package app

import (
	"database/sql"
	"io/ioutil"
	"strings"
)

/// Execute a db command, panic if error occurs
func mustExec(db *sql.DB, command string, args ...interface{}) {
	if _, err := db.Exec(command, args...); err != nil {
		panic(err)
	}
}

/// Execute db commands from script file
func execScriptFile(db *sql.DB, scriptPath string) error {
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

/// Execute db commands from script file, panic if error occurs
func mustExecScriptFile(db *sql.DB, scriptPath string) {
	if err := execScriptFile(db, scriptPath); err != nil {
		panic(err)
	}
}
