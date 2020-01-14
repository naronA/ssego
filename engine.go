package ssego

import (
	"bufio"
	"database/sql"
	"io"
	"os"
	"path/filepath"
)

// 検索エンジンとは？
//   - トークナイザ = ドキュメントの単語をどのようにして区切るか？
//   - インデクサ = ドキュメントからポスティングリストを作成する
//   - ドキュメント管理機 = MySQLのdocumentテーブルモデル
//   - インデクスの保存先ディレクトリパス
type Engine struct {
	tokenizer     *Tokenizer     // トークンの分割方法を決めるトークナイザ
	indexer       *Indexer       // インデクス生成器
	documentStore *DocumentStore // ドキュメント管理機
	indexDir      string         // インデクスファイルを保存するディレクトリ
}

func NewSearchEngine(db *sql.DB) *Engine {
	tokenizer := NewTokenizer()
	indexer := NewIndexer(tokenizer)
	documentStore := NewDocumentStore(db)

	path, ok := os.LookupEnv("INDEX_DIR_PATH")
	if !ok {
		current, _ := os.Getwd()
		path = filepath.Join(current, "_index_data")
	}

	return &Engine{
		tokenizer:     tokenizer,
		indexer:       indexer,
		documentStore: documentStore,
		indexDir:      path,
	}
}

// インデクスにドキュメントを追加する
func (e *Engine) AddDocument(title string, reader io.ReadSeeker) error {
	termCount := e.CountTerm(reader)
	id, err := e.documentStore.save(title, termCount) // タイトルを保存しドキュメントIDを発行する
	if err != nil {
		return err
	}
	reader.Seek(0, 0)
	e.indexer.update(id, reader) // インデクスを更新する
	return nil
}

// ドキュメントをインデクスに追加する処理
func (e *Engine) CountTerm(reader io.Reader) int {
	docLen := 0
	scanner := bufio.NewScanner(reader)
	scanner.Split(e.tokenizer.SplitFunc) // 分割方法の指定
	for scanner.Scan() {
		docLen++
	}
	return docLen
}

func (e *Engine) Flush() error {
	writer := NewIndexWriter(e.indexDir)
	return writer.Flush(e.indexer.index)
}

func (e *Engine) Search(query string, k int, score string) ([]*SearchResult, error) {
	// クエリをトークンに分割
	terms := e.tokenizer.TextToWordSequence(query)

	// 検索を実行
	searcher := NewSearcher(e.indexDir, e.documentStore, score).SearchTopK(terms, k)

	// タイトルを取得
	results := make([]*SearchResult, 0, k)
	for _, result := range searcher.scoreDocs {
		title, err := e.documentStore.fetchTitle(result.docID)
		if err != nil {
			return nil, err
		}
		results = append(results, &SearchResult{
			result.docID, result.score, title,
		})
	}
	return results, nil
}

// 検索結果を格納する構造体
type SearchResult struct {
	DocID DocumentID
	Score float64
	Title string
}
