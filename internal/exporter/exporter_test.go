package exporter

import (
	"os"
	"strings"
	"testing"

	"book-scraper/internal/scraper"
)

func TestSaveJSON(t *testing.T) {
	books := []scraper.Book{{Title: "Book A", Price: 10.0, Rating: 4, InStock: true}}
	path := "test_books.json"
	defer os.Remove(path)

	if err := SaveJSON(books, path); err != nil {
		t.Fatalf("erro ao salvar JSON: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("arquivo JSON não foi criado")
	}
}

func TestSaveCSV(t *testing.T) {
	books := []scraper.Book{{Title: "Book B", Price: 5.5, Rating: 2, InStock: false}}
	path := "test_books.csv"
	defer os.Remove(path)

	if err := SaveCSV(books, path); err != nil {
		t.Fatalf("erro ao salvar CSV: %v", err)
	}
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatal("erro ao ler CSV")
	}
	if !strings.Contains(string(content), "Book B") {
		t.Error("CSV não contém 'Book B'")
	}
}