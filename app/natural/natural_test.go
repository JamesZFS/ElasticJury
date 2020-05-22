package natural

import (
	"testing"
)

func TestParseFullText(t *testing.T) {
	Initialize()
	defer Finalize()
	type args struct {
		text string
	}
	tests := []struct {
		name       string
		args       args
		wantWords  []string
		wantTags   []string
		wantLaws   []string
		wantJudges []string
	}{
		{
			name:       "nil",
			args:       args{""},
			wantWords:  nil,
			wantTags:   nil,
			wantLaws:   nil,
			wantJudges: nil,
		},
		{
			name:       "only words",
			args:       args{"我到清华大学念书"},
			wantWords:  []string{"清华大学", "念书"},
			wantTags:   nil,
			wantLaws:   nil,
			wantJudges: nil,
		},
		{
			name:       "only laws",
			args:       args{"《自强不息》《厚德载物》"},
			wantWords:  nil,
			wantTags:   nil,
			wantLaws:   []string{"《自强不息》", "《厚德载物》"},
			wantJudges: nil,
		},
		{
			name:       "bad laws 1",
			args:       args{"《自强不息》》"},
			wantWords:  nil,
			wantTags:   nil,
			wantLaws:   []string{"《自强不息》"},
			wantJudges: nil,
		},
		{
			name:       "bad laws 2",
			args:       args{"《自《强不息》, 《厚德载物》"},
			wantWords:  nil,
			wantTags:   nil,
			wantLaws:   []string{"《强不息》", "《厚德载物》"},
			wantJudges: nil,
		},
		{
			name:       "bad laws 3",
			args:       args{"《》》《《》《x》"},
			wantWords:  nil,
			wantTags:   nil,
			wantLaws:   []string{"《x》"},
			wantJudges: nil,
		},
		{
			name:       "misc 1",
			args:       args{"我在清华大学学习《x《搜索引擎技术》》"},
			wantWords:  []string{"清华大学", "学习"},
			wantTags:   nil,
			wantLaws:   []string{"《搜索引擎技术》"},
			wantJudges: nil,
		},
		{
			name:       "misc 2",
			args:       args{"我在清华大学学习《x《搜索引擎技术》》这是一门非常重要的课《《法律》》完"},
			wantWords:  []string{"清华大学", "学习", "这是", "一门", "课", "完"},
			wantTags:   nil,
			wantLaws:   []string{"《搜索引擎技术》", "《法律》"},
			wantJudges: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWords, gotTags, gotLaws, gotJudges := ParseFullText(tt.args.text)
			if !Equal(gotWords, tt.wantWords) {
				t.Errorf("ParseFullText() gotWords = %v, want %v", gotWords, tt.wantWords)
			}
			if !Equal(gotTags, tt.wantTags) {
				t.Errorf("ParseFullText() gotTags = %v, want %v", gotTags, tt.wantTags)
			}
			if !Equal(gotLaws, tt.wantLaws) {
				t.Errorf("ParseFullText() gotLaws = %v, want %v", gotLaws, tt.wantLaws)
			}
			if !Equal(gotJudges, tt.wantJudges) {
				t.Errorf("ParseFullText() gotJudges = %v, want %v", gotJudges, tt.wantJudges)
			}
		})
	}
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
