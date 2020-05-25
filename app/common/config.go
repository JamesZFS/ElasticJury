package common

import "time"

const (
	AppDatabaseName        	= "ElasticJury"
	DataSourceName         	= "root:<password>@tcp(cdb-f0b6x25m.cd.tencentcdb.com:10104)/<database>"
	DBConnMaxLifeTime      	= 10 * time.Minute // should be smaller than `select @@global.wait_timeout`
	InitTableScriptPath    	= "database/init-tables.sql"
	SearchFilter		   	= 0.2
	SearchLimit            	= 0 // zero for no limit
	SearchWordLimit		   	= 20
	TipsCount			   	= 10
	StopWordsPath          	= "dicts/stopwords.txt"
	TagDictPath			   	= "dicts/tags.csv"
	JudgeDictPath		   	= "dicts/judges.csv"
	LawsDictPath		  	= "dicts/laws.csv"
)
