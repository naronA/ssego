package ssego

type DocumentStats struct {
	TermCounts map[DocumentID]int
}

func NewDocumentStas() *DocumentStats {
	return &DocumentStats{
		TermCounts: make(map[DocumentID]int),
	}
}
