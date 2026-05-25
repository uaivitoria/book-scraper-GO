package scraper

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

type Book struct {
	Title   string  `json:"title"`
	Price   float64 `json:"price_gbp"`
	Rating  int     `json:"rating"`
	InStock bool    `json:"in_stock"`
}

var RatingMap = map[string]int{
	"One": 1, "Two": 2, "Three": 3, "Four": 4, "Five": 5,
}

const baseURL = "https://books.toscrape.com/catalogue/"
const userAgent = "BookScraper/1.0 (trainee-challenge; educational use)"

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

func getAttribute(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

func hasClass(n *html.Node, class string) bool {
	return strings.Contains(getAttribute(n, "class"), class)
}

func getText(n *html.Node) string {
	if n != nil && n.FirstChild != nil {
		return strings.TrimSpace(n.FirstChild.Data)
	}
	return ""
}

func FindFirst(n *html.Node, match func(*html.Node) bool) *html.Node {
	if match(n) {
		return n
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := FindFirst(c, match); found != nil {
			return found
		}
	}
	return nil
}

func walkNodes(n *html.Node, match func(*html.Node) bool, results *[]*html.Node) {
	if match(n) {
		*results = append(*results, n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		walkNodes(c, match, results)
	}
}

func FindAll(n *html.Node, match func(*html.Node) bool) []*html.Node {
	var results []*html.Node
	walkNodes(n, match, &results)
	return results
}

func ParseBook(article *html.Node) Book {
	h3 := FindFirst(article, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "h3"
	})
	titleNode := FindFirst(h3, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "a"
	})
	title := getAttribute(titleNode, "title")

	priceNode := FindFirst(article, func(n *html.Node) bool {
		return n.Type == html.ElementNode && hasClass(n, "price_color")
	})
	priceText := strings.TrimSpace(getText(priceNode))
	priceText = strings.ReplaceAll(priceText, "£", "")
	priceText = strings.ReplaceAll(priceText, "Â", "")
	priceText = strings.TrimSpace(priceText)
	price, _ := strconv.ParseFloat(priceText, 64)

	ratingNode := FindFirst(article, func(n *html.Node) bool {
		return n.Type == html.ElementNode && hasClass(n, "star-rating")
	})
	ratingClass := getAttribute(ratingNode, "class")
	parts := strings.Split(ratingClass, " ")
	rating := 0
	if len(parts) > 1 {
		rating = RatingMap[parts[1]]
	}

	availNode := FindFirst(article, func(n *html.Node) bool {
		return n.Type == html.ElementNode && hasClass(n, "availability")
	})
	availText := strings.ToLower(strings.TrimSpace(getText(availNode)))
	inStock := strings.Contains(availText, "in stock")

	return Book{Title: title, Price: price, Rating: rating, InStock: inStock}
}

func ScrapeAll() []Book {
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

func scrapePage(url string) ([]Book, string, error) {
	doc, err := fetchPage(url)
	if err != nil {
		return nil, "", err
	}

	articles := FindAll(doc, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "article" && hasClass(n, "product_pod")
	})

	var books []Book
	for _, article := range articles {
		books = append(books, ParseBook(article))
	}

	nextNode := FindFirst(doc, func(n *html.Node) bool {
		return n.Type == html.ElementNode && n.Data == "li" && hasClass(n, "next")
	})
	nextURL := ""
	if nextNode != nil {
		aNode := FindFirst(nextNode, func(n *html.Node) bool {
			return n.Type == html.ElementNode && n.Data == "a"
		})
		if aNode != nil {
			nextURL = baseURL + getAttribute(aNode, "href")
		}
	}

	return books, nextURL, nil
}