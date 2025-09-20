package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const iconDir = `C:\documents\BUREAUBLADICONEN`

func main() {
	url := getURL()
	doc, err := fetchSite(url)
	if err != nil {
		log.Fatalf("Error fetching site: %v", err)
	}
	siteName := getSiteName(&doc, url)
	faviconURL := getFavicon(&doc, url, siteName)
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
		return goquery.Document{}, fmt.Errorf("failed to fetch URL: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		return goquery.Document{}, fmt.Errorf("request failed with status: %s", res.Status)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		return goquery.Document{}, fmt.Errorf("failed to parse HTML: %v", err)
	}

	return *doc, nil
}

func getSiteName(doc *goquery.Document, inputURL string) string {
	siteName := strings.TrimSpace(doc.Find("title").First().Text())
	if siteName == "" {
		parsedURL, _ := url.Parse(inputURL)
		siteName = parsedURL.Host
	}
	fmt.Printf("Found site name: %s\n", siteName)
	return siteName
}

func getFavicon(doc *goquery.Document, inputURL, siteName string) string {
	faviconURL, exists := doc.Find("link[rel='icon'], link[rel='shortcut icon']").Attr("href")
	if !exists {
		parsedURL, _ := url.Parse(inputURL)
		faviconURL = parsedURL.Scheme + "://" + parsedURL.Host + "/favicon.ico"
		fmt.Println("No favicon link found, trying default /favicon.ico")
	} else {
		base, _ := url.Parse(inputURL)
		iconURL, _ := url.Parse(faviconURL)
		faviconURL = base.ResolveReference(iconURL).String()
	}
	fmt.Printf("Found favicon URL: %s\n", faviconURL)

	return faviconURL
}
