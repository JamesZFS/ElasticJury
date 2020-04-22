package app

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"path"
	"reflect"
	"testing"
)

func prologue() *sql.DB {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		panic(err)
	}
	// new test database
	// language=MySQL
	{
		mustExec(db, "DROP DATABASE IF EXISTS ElasticJury_test")
		mustExec(db, "CREATE DATABASE ElasticJury_test DEFAULT CHARACTER SET utf8")
		mustExec(db, "USE ElasticJury_test")
	}
	mustExecScriptFile(db, path.Join("../", initTableScriptPath))
	mustExecScriptFile(db, path.Join("../", initTestDataScriptPath))
	return db
}

func epilogue(db *sql.DB) {
	mustExec(db, "DROP DATABASE ElasticJury_test")
	if err := db.Close(); err != nil {
		panic(err)
	}
}

func Test_searchCaseIdsByWords(t *testing.T) {
	db := prologue()
	defer epilogue(db)
	type args struct {
		words []string

		mode bool
	}
	tests := []struct { // Test cases:
		name    string
		args    args
		want    searchResultSet
		wantErr bool
	}{
		{
			name: "No word",
			args: args{},
			want: searchResultSet{},
		},
		{
			name: "Single word found 1",
			args: args{words: []string{"你好"}},
			want: searchResultSet{
				1: 1.0,
			},
		},
		{
			name: "Single word found 2",
			args: args{words: []string{"某人"}},
			want: searchResultSet{
				1: 0.3,
				2: 0.3,
				3: 0.3,
			},
		},
		{
			name: "Single word not found",
			args: args{words: []string{"啥"}},
			want: searchResultSet{},
		},
		{
			name: "Multi words modeAnd 1",
			args: args{words: []string{"你好", "世界"}, mode: modeAnd},
			want: searchResultSet{
				1: 1.3,
			},
		},
		{
			name: "Multi words modeAnd 2",
			args: args{words: []string{"猥亵", "某人"}, mode: modeAnd},
			want: searchResultSet{
				2: 0.8,
				3: 0.8,
			},
		},
		{
			name: "Multi words modeAnd 3",
			args: args{words: []string{"猥亵", "某人", "what"}, mode: modeAnd},
			want: searchResultSet{},
		},
		{
			name: "Multi words modeOr 1",
			args: args{words: []string{"猥亵", "某人", "what"}, mode: modeOr},
			want: searchResultSet{
				1: 0.3,
				2: 0.8,
				3: 0.8,
			},
		},
		{
			name: "Multi words modeOr 2",
			args: args{words: []string{"世界", "你好"}, mode: modeOr},
			want: searchResultSet{
				1: 1.3,
				2: 0.3,
				3: 0.3,
			},
		},
		{
			name: "Multi words modeOr 3",
			args: args{words: []string{"what", "that"}, mode: modeOr},
			want: searchResultSet{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchCaseIdsByWords(db, tt.args.words, tt.args.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("searchCaseIdsByWords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchCaseIdsByWords() got = %v, want %v", got, tt.want)
			}
		})
	}
}
