package scraper

import (
	"strings"
	"testing"

	"golang.org/x/net/html"
)

func htmlToNode(raw string) *html.Node {
	node, _ := html.Parse(strings.NewReader(raw))
	return node
}

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
	return FindFirst(doc, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "article"
	})
}

func TestParseBookTitle(t *testing.T) {
	article := makeArticle("A Great Book", "£9.99", "Three", "In stock")
	book := ParseBook(article)
	if book.Title != "A Great Book" {
		t.Errorf("esperado 'A Great Book', got '%s'", book.Title)
	}
}

func TestParseBookPrice(t *testing.T) {
	article := makeArticle("Some Book Here", "£9.99", "Three", "In stock")
	book := ParseBook(article)
	if book.Price != 9.99 {
		t.Errorf("esperado 9.99, got %f", book.Price)
	}
}

func TestParseBookRatingThree(t *testing.T) {
	article := makeArticle("Some Book Here", "£9.99", "Three", "In stock")
	book := ParseBook(article)
	if book.Rating != 3 {
		t.Errorf("esperado 3, got %d", book.Rating)
	}
}

func TestParseBookRatingFive(t *testing.T) {
	article := makeArticle("Some Book Here", "£9.99", "Five", "In stock")
	book := ParseBook(article)
	if book.Rating != 5 {
		t.Errorf("esperado 5, got %d", book.Rating)
	}
}

func TestParseBookInStock(t *testing.T) {
	article := makeArticle("Some Book Here", "£9.99", "Three", "In stock")
	book := ParseBook(article)
	if !book.InStock {
		t.Error("esperado InStock = true")
	}
}

func TestParseBookOutOfStock(t *testing.T) {
	article := makeArticle("Some Book Here", "£9.99", "Three", "Out of stock")
	book := ParseBook(article)
	if book.InStock {
		t.Error("esperado InStock = false")
	}
}

func TestRatingMapComplete(t *testing.T) {
	expected := map[int]bool{1: false, 2: false, 3: false, 4: false, 5: false}
	for _, v := range RatingMap {
		expected[v] = true
	}
	for k, found := range expected {
		if !found {
			t.Errorf("rating %d não encontrado no RatingMap", k)
		}
	}
}