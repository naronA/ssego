package ssego

import (
	"reflect"
	"testing"
)

func TestSearchTopK(t *testing.T) {
	s := NewSearcher("testdata/index", nil, "TFIDF")      // searcherの初期化
	actual := s.SearchTopK([]string{"quarrel", "sir"}, 1) // 検索の実行

	expected := &TopDocs{2, []*ScoreDoc{{2, 1.9657842846620868}}}

	for !reflect.DeepEqual(actual, expected) {
		t.Fatalf("got:%v\nexpected:%v\n", actual, expected)
	}
}
