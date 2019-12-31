package ssego

import (
	"bufio"
	"io"
)

type Indexer struct {
	index     *TransposeIndex
	tokenizer *Tokenizer
}

func NewIndexer(tokenizer *Tokenizer) *Indexer {
	return &Indexer{
		index:     NewIndex(),
		tokenizer: tokenizer,
	}
}

// ドキュメントをインデクスに追加する処理
func (idxr *Indexer) update(docID DocumentID, reader io.Reader) {

	// bufio.Scannerを使用することでファイルや標準入力などからデータを少しずつ読み込むことができる。
	scanner := bufio.NewScanner(reader)
	scanner.Split(idxr.tokenizer.SplitFunc) // 分割方法の指定
	var position int

	for scanner.Scan() {
		term := scanner.Text() // 用語ごとに読み込み

		// ポスティングリストの更新
		if postingsList, ok := idxr.index.Dictionary[term]; !ok {
			// termをキーとするポスティングリストが存在しない場合
			idxr.index.Dictionary[term] = NewPostingsList(NewPosting(docID, position))
		} else {
			// ポスティングリストがすでに存在する場合は追加
			postingsList.Add(NewPosting(docID, position))
		}
		position++
	}

	idxr.index.TotalDocsCount++

}
