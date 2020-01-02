package ssego

import (
	"math"
	"sort"
)

// 検索処理を担う構造体Searcher
type Searcher struct {
	indexReader   *IndexReader // インデクス読み取り器
	cursors       []*Cursor    // ポスティングリストのポインタ配列
	documentStats *DocumentStats
}

func NewSearcher(path string, docStats *DocumentStats) *Searcher {
	return &Searcher{indexReader: NewIndexReader(path), documentStats: docStats}
}

// 検索を実行し、スコアが高い順にK件結果を返す
func (s *Searcher) SearchTopK(query []string, k int) *TopDocs {
	// マッチするドキュメントを抽出しスコアを計算する
	results := s.search(query)

	// 結果をスコアの降順でソートする
	sort.Slice(results, func(i, j int) bool {
		return results[i].score > results[j].score
	})

	total := len(results)
	if len(results) > k {
		results = results[:k] // 上位k件のみ取得
	}

	return &TopDocs{
		totalHits: total,
		scoreDocs: results,
	}
}

func (s *Searcher) search(query []string) []*ScoreDoc {
	// カーソルの取得
	// クエリに含まれる用語のポスティングリストが一つも存在しない場合、0件で終了する
	if s.openCursors(query) == 0 {
		return []*ScoreDoc{}
	}

	// 一番短いポスティングリストを参照するカーソルを洗濯
	c := s.cursors[0]
	cursors := s.cursors[1:]

	// 結果を格納する構造体の初期化
	docs := make([]*ScoreDoc, 0)

	scorer := &TFIDFScore{indexReader: s.indexReader, cursors: s.cursors}
	// 最も短いポスティングリストをたどり終えるまで繰り返す
	for !c.Empty() {
		var nextDocID DocumentID
		// その他のカーソルをcのdocID以上になるまですすめる
		for _, cursor := range cursors {
			if cursor.NextDoc(c.DocID()); cursor.Empty() {
				return docs
			}
			// docIDが一致しなければ
			if cursor.DocID() != c.DocID() {
				nextDocID = cursor.DocID()
				break
			}
		}

		if nextDocID > 0 {
			// nextDocID以上になるまで読みすすめる
			if c.NextDoc(nextDocID); c.Empty() {
				return docs
			}
		} else {
			// 結果を格納
			docs = append(docs, &ScoreDoc{
				docID: c.DocID(),
				score: scorer.CalcScore(),
			})
			c.Next()
		}
	}
	return docs

}

func (s *Searcher) openCursors(query []string) int {
	// ポスティングリストを取得
	postings := s.indexReader.postingsLists(query)
	if len(postings) == 0 {
		return 0
	}

	// ポスティングリストの短い順にソート
	sort.Slice(postings, func(i, j int) bool {
		return postings[i].Len() < postings[j].Len()
	})

	// 各ポスティングリストに対するcursorの取得
	cursors := make([]*Cursor, len(postings))
	for i, postingList := range postings {
		cursors[i] = postingList.OpenCursor()
	}

	s.cursors = cursors
	return len(cursors)
}

type Scorer interface {
	CalcScore() float64
}

type TFIDFScore struct {
	indexReader *IndexReader // インデクス読み取り器
	cursors     []*Cursor    // ポスティングリストのポインタ配列
}

func (t TFIDFScore) CalcScore() float64 {
	var score float64
	for i := 0; i < len(t.cursors); i++ {
		termFreq := t.cursors[i].Posting().TermFrequency
		docCount := t.cursors[i].postingsList.Len()
		totalDocCount := t.indexReader.totalDocCount()
		score += calcTF(termFreq) * calcIDF(totalDocCount, docCount)
	}
	return score

}

type BM25Score struct {
	indexReader *IndexReader // インデクス読み取り器
	cursors     []*Cursor    // ポスティングリストのポインタ配列
}

func (b BM25Score) CalcScore() float64 {
	var score float64
	// s.documentStats[]
	for i := 0; i < len(b.cursors); i++ {
		// docID := s.cursors[i].DocID()
		// termCount := s.documentStats.TermCounts[docID]
		termFreq := b.cursors[i].Posting().TermFrequency
		docCount := b.cursors[i].postingsList.Len()
		totalDocCount := b.indexReader.totalDocCount()
		score += calcTF(termFreq) * calcIDF(totalDocCount, docCount)
	}
	return score

}

func calcTF(termCount int) float64 {
	if termCount <= 0 {
		return 0
	}

	return math.Log2(float64(termCount)) + 1
}

// Inverse Document Frequency
// 全ドキュメント数 N と 用語が含まれているドキュメント数 dfを用いてIDFを計算する
func calcIDF(N, df int) float64 {
	return math.Log2(float64(N) / float64(df))
}
