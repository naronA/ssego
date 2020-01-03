package ssego

import (
	"database/sql"
	"log"
)

type DocumentStore struct {
	db *sql.DB
}

func NewDocumentStore(db *sql.DB) *DocumentStore {
	return &DocumentStore{db: db}
}

func (ds *DocumentStore) save(title string, termCount int) (DocumentID, error) {
	query := "INSERT INTO documents (document_title, document_terms) VALUES (?, ?)"
	result, err := ds.db.Exec(query, title, termCount)
	if err != nil {
		log.Fatal(err)
	}
	id, err := result.LastInsertId()
	return DocumentID(id), err
}

func (ds *DocumentStore) fetchTitle(docID DocumentID) (string, error) {
	query := "SELECT document_title FROM documents WHERE document_id = ?"
	row := ds.db.QueryRow(query, docID)
	var title string
	err := row.Scan(&title)
	return title, err
}

func (ds *DocumentStore) fetchTermCount(docID DocumentID) (int, error) {
	query := "SELECT document_terms FROM documents WHERE document_id = ?"
	row := ds.db.QueryRow(query, docID)
	var termCount int
	err := row.Scan(&termCount)
	return termCount, err
}
