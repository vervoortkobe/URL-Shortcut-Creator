package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const iconDir = `C:\documents\BUREAUBLADICONEN`

func main() {
	url := getURL()
	fetchSite(url)
}

func getURL() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter the URL to create a shortcut for: ")
	inputURL, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("Failed to read input: %v", err)
	}
	inputURL = strings.TrimSpace(inputURL)

	if !strings.HasPrefix(inputURL, "http://") && !strings.HasPrefix(inputURL, "https://") {
		inputURL = "https://" + inputURL
	}

	fmt.Printf("Processing URL: %s\n", inputURL)

	return inputURL
}

func fetchSite(inputURL string) (goquery.Document, error) {
	res, err := http.Get(inputURL)
	if err != nil {
		return goquery.Document{}, fmt.Errorf("Failed to fetch URL: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return goquery.Document{}, fmt.Errorf("Request failed with status: %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return goquery.Document{}, fmt.Errorf("Failed to parse HTML: %v", err)
	}

	return *doc, nil
}
