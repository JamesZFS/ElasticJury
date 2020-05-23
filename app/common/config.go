package common

import "time"

const (
	AppDatabaseName        = "ElasticJury"
	TestDatabaseName       = "Test"
	DataSourceName         = "root:<password>@tcp(cdb-f0b6x25m.cd.tencentcdb.com:10104)/<database>"
	DBConnMaxLifeTime      = 10 * time.Minute // should be smaller than `select @@global.wait_timeout`
	InitTableScriptPath    = "database/init-tables.sql"
	InitTestDataScriptPath = "database/init-test-data.sql"
	StopWordsPath          = "preprocessor/stopwords.txt"
	WordSearchLimit        = 200
)
