package exporter

import (
	"encoding/csv"
	"encoding/json"
	"os"
	"strconv"

	"book-scraper/internal/scraper"
)

func SaveJSON(books []scraper.Book, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(books)
}

func SaveCSV(books []scraper.Book, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	writer.Write([]string{"title", "price_gbp", "rating", "in_stock"})
	for _, b := range books {
		writer.Write([]string{
			b.Title,
			strconv.FormatFloat(b.Price, 'f', 2, 64),
			strconv.Itoa(b.Rating),
			strconv.FormatBool(b.InStock),
		})
	}
	return nil
}