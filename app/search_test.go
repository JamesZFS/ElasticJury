package app

import (
	. "ElasticJury/app/common"
	"fmt"
	"reflect"
	"testing"
)

func Test_searchCaseIdsByWord(t *testing.T) {
	db := dbPrologue()
	defer dbEpilogue(db)
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
			got, err := db.searchCaseIdsByWord(tt.args.words, tt.args.mode)
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
	db := dbPrologue()
	defer dbEpilogue(db)
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
			got, err := db.searchCaseIdsByTag(tt.args.tags, tt.args.mode)
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
		s searchResultSet
		t searchResultSet
	}
	tests := []struct {
		name string
		args args
		want searchResultSet
	}{
		{
			name: "Merge empty sets",
			args: args{
				s: searchResultSet{},
				t: searchResultSet{},
			},
			want: searchResultSet{},
		},
		{
			name: "Merge not empty 1",
			args: args{
				s: searchResultSet{
					1: 0.5,
					2: 0.8,
					3: 0.1,
				},
				t: searchResultSet{
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
				s: searchResultSet{
					2: 0.1,
					3: 0.0,
				},
				t: searchResultSet{
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
				s: searchResultSet{
					1: 0.5,
					2: 0.8,
					3: 0.1,
				},
				t: searchResultSet{
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
				s: searchResultSet{
					2: 0.1,
					3: 0.8,
				},
				t: searchResultSet{
					1: 0.1,
					4: 0.1,
				},
			},
			want: searchResultSet{},
		},
		{
			name: "Merge with nil 1",
			args: args{
				s: nil,
				t: searchResultSet{
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
				s: nil,
				t: nil,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.s.merge(tt.args.t); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("merge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func dbPrologue() database {
	password := GetEnvVar("PASSWORD", "")
	dbRoot, err := newDatabase("", password)
	if err != nil {
		panic(err)
	}
	// drop old test database
	dbRoot.mustExec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", TestDatabaseName))
	dbRoot.mustExec(fmt.Sprintf("CREATE DATABASE %s DEFAULT CHARACTER SET utf8", TestDatabaseName))
	db, err := newDatabase(TestDatabaseName, password)
	if err != nil {
		panic(err)
	}
	db.mustExecScriptFile(InitTableScriptPath)
	db.mustExecScriptFile(InitTestDataScriptPath)
	return db
}

func dbEpilogue(db database) {
	db.mustExec(fmt.Sprintf("DROP DATABASE %s", TestDatabaseName))
	if err := db.Close(); err != nil {
		panic(err)
	}
}
