package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

// Book representa os dados de um livro coletado
type Book struct {
	Title   string  `json:"title"`
	Price   float64 `json:"price_gbp"`
	Rating  int     `json:"rating"`
	InStock bool    `json:"in_stock"`
}

// mapa de rating texto → número
var ratingMap = map[string]int{
	"One": 1, "Two": 2, "Three": 3, "Four": 4, "Five": 5,
}

const baseURL = "https://books.toscrape.com/catalogue/"
const userAgent = "BookScraper/1.0 (trainee-challenge; educational use)"

// fetchPage busca o HTML de uma URL e retorna o nó raiz
func fetchPage(url string) (*html.Node, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return html.Parse(resp.Body)
}

// getAttribute retorna o valor de um atributo HTML
func getAttribute(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// hasClass verifica se um nó tem uma classe CSS específica
func hasClass(n *html.Node, class string) bool {
	return strings.Contains(getAttribute(n, "class"), class)
}

// getText retorna o texto interno de um nó
func getText(n *html.Node) string {
	if n != nil && n.FirstChild != nil {
		return strings.TrimSpace(n.FirstChild.Data)
	}
	return ""
}

// findFirst busca recursivamente o primeiro nó que satisfaz a condição
func findFirst(n *html.Node, match func(*html.Node) bool) *html.Node {
	if match(n) {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := findFirst(c, match); found != nil {
			return found
		}
	}
	return nil
}

// findAll busca recursivamente todos os nós que satisfazem a condição
func findAll(n *html.Node, match func(*html.Node) bool) []*html.Node {
	var results []*html.Node
	if match(n) {
		results = append(results, n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		results = append(results, findAll(c, match)...)
	}
	return results
}

// parseBook extrai os dados de um nó <article class="product_pod">
func parseBook(article *html.Node) Book {
	// Título
	h3 := findFirst(article, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "h3"
	})
	titleNode := findFirst(h3, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "a"
	})
	title := getAttribute(titleNode, "title")

	// Preço
	priceNode := findFirst(article, func(n *html.Node) bool {
		return n.Type == html.ElementNode && hasClass(n, "price_color")
	})
	priceText := strings.TrimSpace(getText(priceNode))
	priceText = strings.ReplaceAll(priceText, "£", "")
	priceText = strings.ReplaceAll(priceText, "Â", "")
	priceText = strings.TrimSpace(priceText)
	price, _ := strconv.ParseFloat(priceText, 64)

	// Rating
	ratingNode := findFirst(article, func(n *html.Node) bool {
		return n.Type == html.ElementNode && hasClass(n, "star-rating")
	})
	ratingClass := getAttribute(ratingNode, "class")
	parts := strings.Split(ratingClass, " ")
	rating := 0
	if len(parts) > 1 {
		rating = ratingMap[parts[1]]
	}

	// Disponibilidade
	availNode := findFirst(article, func(n *html.Node) bool {
		return n.Type == html.ElementNode && hasClass(n, "availability")
	})
	availText := strings.ToLower(strings.TrimSpace(getText(availNode)))
	inStock := strings.Contains(availText, "in stock")

	return Book{Title: title, Price: price, Rating: rating, InStock: inStock}
}

// scrapePage coleta os livros de uma página e retorna a URL da próxima
func scrapePage(url string) ([]Book, string, error) {
	doc, err := fetchPage(url)
	if err != nil {
		return nil, "", err
	}

	// Encontra todos os <article class="product_pod">
	articles := findAll(doc, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "article" && hasClass(n, "product_pod")
	})

	var books []Book
	for _, article := range articles {
		books = append(books, parseBook(article))
	}

	// Próxima página
	nextNode := findFirst(doc, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "li" && hasClass(n, "next")
	})
	nextURL := ""
	if nextNode != nil {
		aNode := findFirst(nextNode, func(n *html.Node) bool {
			return n.Type == html.ElementNode && n.Data == "a"
		})
		if aNode != nil {
			nextURL = baseURL + getAttribute(aNode, "href")
		}
	}

	return books, nextURL, nil
}

// scrapeAll percorre todas as páginas e retorna todos os livros
func scrapeAll() []Book {
	var allBooks []Book
	url := baseURL + "page-1.html"
	page := 1

	for url != "" {
		fmt.Printf("Coletando página %d...\n", page)
		books, nextURL, err := scrapePage(url)
		if err != nil {
			log.Printf("Erro na página %d: %v", page, err)
			break
		}
		allBooks = append(allBooks, books...)
		url = nextURL
		page++
		if url != "" {
			time.Sleep(500 * time.Millisecond)
		}
	}

	fmt.Printf("Total coletado: %d livros.\n", len(allBooks))
	return allBooks
}

// saveJSON salva os livros em formato JSON
func saveJSON(books []Book, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	return encoder.Encode(books)
}

// saveCSV salva os livros em formato CSV
func saveCSV(books []Book, path string) error {
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

func main() {
	books := scrapeAll()

	if err := saveJSON(books, "books.json"); err != nil {
		log.Fatal("Erro ao salvar JSON:", err)
	}
	fmt.Println("JSON salvo em: books.json")

	if err := saveCSV(books, "books.csv"); err != nil {
		log.Fatal("Erro ao salvar CSV:", err)
	}
	fmt.Println("CSV salvo em: books.csv")
}