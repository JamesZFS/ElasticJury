package natural

import (
	. "ElasticJury/app/common"
	"reflect"
	"testing"
)

func subSet(words1 []string, words2 []string) bool {
	for _, w := range words1 {
		if IndexOfStr(words2, w) < 0 { // not found
			return false
		}
	}
	return true
}

// Equal tells whether a and b contain the same elements.
// A nil argument is equivalent to an empty slice.
func Equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestPreprocessWord(t *testing.T) {
	type args struct {
		word string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "1",
			args: args{"审判"},
			want: "审判",
		},
		{
			name: "2",
			args: args{"  \t审判 "},
			want: "审判",
		},
		{
			name: "3",
			args: args{"  \t审 判 "},
			want: "审 判",
		},
		{
			name: "escape",
			args: args{"  \t'《\"审 `判\\》 "},
			want: "《审 判》",
		},
		{
			name: "nil",
			args: args{"  \t'\\\" "},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PreprocessWord(tt.args.word); got != tt.want {
				t.Errorf("PreprocessWord() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPreprocessWords(t *testing.T) {
	type args struct {
		words []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "nil",
			args: args{nil},
			want: nil,
		},
		{
			name: "1",
			args: args{[]string{"123", "321"}},
			want: []string{"123", "321"},
		},
		{
			name: "2",
			args: args{[]string{"1`2\\3", "3\"2'1"}},
			want: []string{"123", "321"},
		},
		{
			name: "3",
			args: args{[]string{"", "\ttag", "\t'`", "《案件》 "}},
			want: []string{"tag", "《案件》"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PreprocessWords(tt.args.words); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PreprocessWords() = %v, want %v", got, tt.want)
			}
		})
	}
}
