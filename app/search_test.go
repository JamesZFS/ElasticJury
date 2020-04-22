package app

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"path"
	"reflect"
	"testing"
)

func Test_searchCaseIdsByWord(t *testing.T) {
	db := prologue()
	defer epilogue(db)
	type args struct {
		words []string
		mode  bool
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
			name: "Multi tags modeAnd 1",
			args: args{words: []string{"你好", "世界"}, mode: modeAnd},
			want: searchResultSet{
				1: 1.3,
			},
		},
		{
			name: "Multi tags modeAnd 2",
			args: args{words: []string{"猥亵", "某人"}, mode: modeAnd},
			want: searchResultSet{
				2: 0.8,
				3: 0.8,
			},
		},
		{
			name: "Multi tags modeAnd 3",
			args: args{words: []string{"猥亵", "某人", "what"}, mode: modeAnd},
			want: searchResultSet{},
		},
		{
			name: "Multi tags modeOr 1",
			args: args{words: []string{"猥亵", "某人", "what"}, mode: modeOr},
			want: searchResultSet{
				1: 0.3,
				2: 0.8,
				3: 0.8,
			},
		},
		{
			name: "Multi tags modeOr 2",
			args: args{words: []string{"世界", "你好"}, mode: modeOr},
			want: searchResultSet{
				1: 1.3,
				2: 0.3,
				3: 0.3,
			},
		},
		{
			name: "Multi tags modeOr 3",
			args: args{words: []string{"what", "that"}, mode: modeOr},
			want: searchResultSet{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchCaseIdsByWord(db, tt.args.words, tt.args.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("searchCaseIdsByWord() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchCaseIdsByWord() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_searchCaseIdsByTag(t *testing.T) {
	db := prologue()
	defer epilogue(db)
	type args struct {
		tags []string
		mode bool
	}
	tests := []struct { // Test cases:
		name    string
		args    args
		want    searchResultSet
		wantErr bool
	}{
		{
			name: "No tag",
			args: args{},
			want: searchResultSet{},
		},
		{
			name: "Single tag found",
			args: args{tags: []string{"诈骗"}},
			want: searchResultSet{
				1: 1.0,
			},
		},
		{
			name: "Single tag not found",
			args: args{tags: []string{"啥"}},
			want: searchResultSet{},
		},
		{
			name: "Multi tags modeAnd",
			args: args{tags: []string{"诈骗", "猥亵"}, mode: modeAnd},
			want: searchResultSet{},
		},
		{
			name: "Multi tags modeOr",
			args: args{tags: []string{"诈骗", "猥亵"}, mode: modeOr},
			want: searchResultSet{
				1: 1.0,
				2: 0.6,
				3: 0.6,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := searchCaseIdsByTag(db, tt.args.tags, tt.args.mode)
			if (err != nil) != tt.wantErr {
				t.Errorf("searchCaseIdsByTag() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("searchCaseIdsByTag() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mergeSearchResult(t *testing.T) {
	type args struct {
		set1 searchResultSet
		set2 searchResultSet
	}
	tests := []struct {
		name string
		args args
		want searchResultSet
	}{
		{
			name: "Merge empty sets",
			args: args{
				set1: searchResultSet{},
				set2: searchResultSet{},
			},
			want: searchResultSet{},
		},
		{
			name: "Merge not empty 1",
			args: args{
				set1: searchResultSet{
					1: 0.5,
					2: 0.8,
					3: 0.1,
				},
				set2: searchResultSet{
					1: 0.1,
					2: 0.2,
					3: 0.1,
				},
			},
			want: searchResultSet{
				1: 0.6,
				2: 1.0,
				3: 0.2,
			},
		},
		{
			name: "Merge not empty 2",
			args: args{
				set1: searchResultSet{
					2: 0.1,
					3: 0.0,
				},
				set2: searchResultSet{
					1: 0.1,
					2: 0.1,
				},
			},
			want: searchResultSet{
				2: 0.2,
			},
		},
		{
			name: "Merge not empty 3",
			args: args{
				set1: searchResultSet{
					1: 0.5,
					2: 0.8,
					3: 0.1,
				},
				set2: searchResultSet{
					1: 0.1,
					3: 0.1,
				},
			},
			want: searchResultSet{
				1: 0.6,
				3: 0.2,
			},
		},
		{
			name: "Merge empty",
			args: args{
				set1: searchResultSet{
					2: 0.1,
					3: 0.8,
				},
				set2: searchResultSet{
					1: 0.1,
					4: 0.1,
				},
			},
			want: searchResultSet{},
		},
		{
			name: "Merge with nil 1",
			args: args{
				set1: nil,
				set2: searchResultSet{
					1: 0.4,
					2: 0.2,
				},
			},
			want: searchResultSet{
				1: 0.4,
				2: 0.2,
			},
		},
		{
			name: "Merge with nil 2",
			args: args{
				set1: nil,
				set2: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mergeSearchResult(tt.args.set1, tt.args.set2); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("mergeSearchResult() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
