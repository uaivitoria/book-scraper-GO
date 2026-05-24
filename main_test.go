package main

import (
	"os"
	"strings"
	"testing"

	"golang.org/x/net/html"
)

// htmlToNode converte uma string HTML em um nó para os testes
func htmlToNode(rawHTML string) *html.Node {
	node, _ := html.Parse(strings.NewReader(rawHTML))
	return node
}

// makeArticle cria um <article> HTML falso para simular um livro
func makeArticle(title, price, rating, availability string) *html.Node {
	raw := `<article class="product_pod">
		<h3><a href="#" title="` + title + `">` + title[:10] + `</a></h3>
		<p class="star-rating ` + rating + `"></p>
		<div class="product_price">
			<p class="price_color">` + price + `</p>
			<p class="availability">` + availability + `</p>
		</div>
	</article>`
	doc := htmlToNode(raw)
	return findFirst(doc, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "article"
	})
}

// --- Testes do parseBook ---

func TestParseBookTitle(t *testing.T) {
	article := makeArticle("A Great Book", "£9.99", "Three", "In stock")
	book := parseBook(article)
	if book.Title != "A Great Book" {
		t.Errorf("esperado 'A Great Book', got '%s'", book.Title)
	}
}

func TestParseBookPrice(t *testing.T) {
	article := makeArticle("Some Book Here", "£9.99", "Three", "In stock")
	book := parseBook(article)
	if book.Price != 9.99 {
		t.Errorf("esperado 9.99, got %f", book.Price)
	}
}

func TestParseBookRatingThree(t *testing.T) {
	article := makeArticle("Some Book Here", "£9.99", "Three", "In stock")
	book := parseBook(article)
	if book.Rating != 3 {
		t.Errorf("esperado 3, got %d", book.Rating)
	}
}

func TestParseBookRatingFive(t *testing.T) {
	article := makeArticle("Some Book Here", "£9.99", "Five", "In stock")
	book := parseBook(article)
	if book.Rating != 5 {
		t.Errorf("esperado 5, got %d", book.Rating)
	}
}

func TestParseBookInStock(t *testing.T) {
	article := makeArticle("Some Book Here", "£9.99", "Three", "In stock")
	book := parseBook(article)
	if !book.InStock {
		t.Error("esperado InStock = true")
	}
}

func TestParseBookOutOfStock(t *testing.T) {
	article := makeArticle("Some Book Here", "£9.99", "Three", "Out of stock")
	book := parseBook(article)
	if book.InStock {
		t.Error("esperado InStock = false")
	}
}

func TestRatingMapComplete(t *testing.T) {
	expected := map[int]bool{1: false, 2: false, 3: false, 4: false, 5: false}
	for _, v := range ratingMap {
		expected[v] = true
	}
	for k, found := range expected {
		if !found {
			t.Errorf("rating %d não encontrado no ratingMap", k)
		}
	}
}

func TestSaveJSON(t *testing.T) {
	books := []Book{{Title: "Book A", Price: 10.0, Rating: 4, InStock: true}}
	path := "test_books.json"
	defer os.Remove(path)

	if err := saveJSON(books, path); err != nil {
		t.Fatalf("erro ao salvar JSON: %v", err)
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("arquivo JSON não foi criado")
	}
}

func TestSaveCSV(t *testing.T) {
	books := []Book{{Title: "Book B", Price: 5.5, Rating: 2, InStock: false}}
	path := "test_books.csv"
	defer os.Remove(path)

	if err := saveCSV(books, path); err != nil {
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