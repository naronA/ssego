package ssego

import (
	"container/list"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

type DocumentID int64

// 辞書
//  |
//  +- 単語: PostingsList
//              |
//              +- Posting
//              |    |
//              |    |- 用語が出てきたドキュメントDocID
//              |    |- Positions このドキュメント内の用語の出現いち位置
//              |    +- TermFrequency ドキュメント内の用語の出現回数
//              |
//              +- Posting
//              |    |
//              |    |- 用語が出てきたドキュメントDocID
//              |    |- Positions このドキュメント内の用語の出現いち位置
//              |    +- TermFrequency ドキュメント内の用語の出現回数
//              |
//              +- Posting
//              |    |
//              |    |- 用語が出てきたドキュメントDocID
//              |    |- Positions このドキュメント内の用語の出現いち位置
//              |    +- TermFrequency ドキュメント内の用語の出現回数
//              |
//              +- Posting
//                   |
//                   |- DocID 用語が出てきたドキュメントDocID
//                   |- Positions このドキュメント内の用語の出現いち位置
//                   +- TermFrequency ドキュメント内の用語の出現回数
// 転置インデックス
type TransposeIndex struct {
	Dictionary     map[string]PostingsList // 辞書
	TotalDocsCount int                     // ドキュメントの総数
}

// NewIndex create a new index
func NewIndex() *TransposeIndex {
	dict := make(map[string]PostingsList) // 辞書

	return &TransposeIndex{
		Dictionary:     dict,
		TotalDocsCount: 0,
	}
}

func (idx *TransposeIndex) String() string {
	var padding int

	keys := make([]string, 0, len(idx.Dictionary))

	for k := range idx.Dictionary {
		l := utf8.RuneCountInString(k)
		if padding < l {
			padding = l
		}

		keys = append(keys, k)
	}

	sort.Strings(keys)
	strs := make([]string, len(keys))
	format := "  [%-]" + strconv.Itoa(padding) + "s] -> %s"

	for i, k := range keys {
		if postingList, ok := idx.Dictionary[k]; ok {
			strs[i] = fmt.Sprintf(format, k, postingList.String())
		}
	}

	return fmt.Sprintf("total documents : %v\ndictionary:\n%v\n", idx.TotalDocsCount, strings.Join(strs, "\n"))
}

type PostingsList struct {
	*list.List
}

// PostingをあつめたPostingList
func NewPostingsList(postings ...*Posting) PostingsList {
	l := list.New()

	for _, posting := range postings {
		l.PushBack(posting)
	}

	return PostingsList{l}
}

func (pl PostingsList) add(p *Posting) {
	pl.PushBack(p)
}

func (pl PostingsList) last() *Posting {
	e := pl.List.Back()

	if e == nil {
		return nil
	}

	return e.Value.(*Posting)
}

func (pl PostingsList) Add(new *Posting) {
	last := pl.last()

	if last == nil || last.DocID != new.DocID {
		pl.add(new)
		return
	}

	last.Positions = append(last.Positions, new.Positions...)
	last.TermFrequency++
}

func (pl PostingsList) String() string {
	str := make([]string, 0, pl.Len())
	for e := pl.Front(); e != nil; e.Next() {
		str = append(str, e.Value.(*Posting).String())
	}

	return strings.Join(str, "=>")
}

func (pl PostingsList) MarshalJSON() ([]byte, error) {
	postings := make([]*Posting, 0, pl.Len())

	for e := pl.Front(); e != nil; e = e.Next() {
		postings = append(postings, e.Value.(*Posting))
	}

	return json.Marshal(postings)
}

func (pl *PostingsList) UnmarshalJSON(b []byte) error {
	var postings []*Posting
	if err := json.Unmarshal(b, &postings); err != nil {
		return err
	}

	pl.List = list.New()

	for _, posting := range postings {
		pl.add(posting)
	}

	return nil
}

// 用語が含まれているDocID, ドキュメント内の位置, 出現回数をまとめた構造体
type Posting struct {
	DocID         DocumentID // 単語が含まれているドキュメントのID
	Positions     []int      // 単語の位置
	TermFrequency int        // 単語の出現回数
}

func NewPosting(docID DocumentID, positions ...int) *Posting {
	return &Posting{docID, positions, len(positions)}
}

func (p Posting) String() string {
	return fmt.Sprintf("(%v, %v, %v)", p.DocID, p.Positions, p.TermFrequency)
}
