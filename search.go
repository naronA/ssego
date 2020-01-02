package ssego

import (
	"fmt"
)

// SearchTopKの検索結果を保持する
type TopDocs struct {
	totalHits int         // ヒット件数
	scoreDocs []*ScoreDoc // 検索結果
}

func (t *TopDocs) String() string {
	return fmt.Sprintf("\ntotal hits: %v\nresults: %v\n", t.totalHits, t.scoreDocs)
}

// ドキュメントIDそのドキュメントのスコアを保持する
type ScoreDoc struct {
	docID DocumentID
	score float64
}

func (d ScoreDoc) String() string {
	return fmt.Sprintf("docID: %v, Score: %v", d.docID, d.score)
}
