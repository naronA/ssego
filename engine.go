package ssego

import (
	"database/sql"
	"io"
	"os"
	"path/filepath"
)

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
func (e *Engine) AddDocument(title string, reader io.Reader) error {
	id, err := e.documentStore.save(title) // タイトルを保存しドキュメントIDを発行する
	if err != nil {
		return err
	}
	e.indexer.update(id, reader) // インデクスを更新する
	return nil
}

func (e *Engine) Flush() error {
	writer := NewIndexWriter(e.indexDir)
	return writer.Flush(e.indexer.index)
}
