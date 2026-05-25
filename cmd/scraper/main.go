package main

import (
	"fmt"
	"log"

	"book-scraper/internal/exporter"
	"book-scraper/internal/scraper"
)

func main() {
	books := scraper.ScrapeAll()

	if err := exporter.SaveJSON(books, "books.json"); err != nil {
		log.Fatal("Erro ao salvar JSON:", err)
	}
	fmt.Println("JSON salvo em: books.json")

	if err := exporter.SaveCSV(books, "books.csv"); err != nil {
		log.Fatal("Erro ao salvar CSV:", err)
	}
	fmt.Println("CSV salvo em: books.csv")
}