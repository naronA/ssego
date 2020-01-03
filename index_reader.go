package ssego

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
)

type IndexReader struct {
	indexDir      string                   // インデクスファイルが保存されているディレクトリのパス
	postingsCache map[string]*PostingsList // 読み込んだポスティングリストをキャッシュするフィールド
	docCountCache int                      // インデクスされたドキュメント数をキャッシュするフィールド
}

func NewIndexReader(path string) *IndexReader {
	cache := make(map[string]*PostingsList)
	return &IndexReader{path, cache, -1}
}

// 複数の検索キーワードtermsにマッチするpostingsを読み込んで[]*PostingsListを作成する
func (r *IndexReader) postingsLists(terms []string) []*PostingsList {
	// 複数のキーワードtermsなのでpostingsList"S"になる
	postingsLists := make([]*PostingsList, 0, len(terms))
	for _, term := range terms {
		if postings := r.postings(term); postings != nil {
			postingsLists = append(postingsLists, postings)
		}
	}
	return postingsLists
}

func (r *IndexReader) postings(term string) *PostingsList {
	// すでに取得図意味であればキャッシュを返す
	if postingsList, ok := r.postingsCache[term]; ok {
		return postingsList
	}

	// インデクスファイルの取得
	filename := filepath.Join(r.indexDir, term)
	file, err := os.Open(filename)
	if err != nil {
		return nil
	}
	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil
	}
	var postingsList PostingsList
	err = json.Unmarshal(bytes, &postingsList)
	if err != nil {
		return nil
	}

	// キャッシュの更新
	r.postingsCache[term] = &postingsList
	return &postingsList
}

func (r *IndexReader) totalDocCount() int {
	// すでに取得済みであればキャッシュを返す
	if r.docCountCache > 0 {
		return r.docCountCache
	}
	filename := filepath.Join(r.indexDir, "_0.dc")
	file, err := os.Open(filename)
	if err != nil {
		// 読み込みに失敗したら0件とする
		return 0
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return 0
	}
	count, err := strconv.Atoi(string(bytes))
	if err != nil {
		return 0
	}
	r.docCountCache = count
	return count
}
