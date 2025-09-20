package main

import (
	"bufio"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/sergeymakinen/go-ico"

	"github.com/PuerkitoBio/goquery"
)

var iconDir string

func main() {
	initIconDir()

	url := getURL()
	doc, err := fetchSite(url)
	if err != nil {
		log.Fatalf("Error fetching site: %v", err)
	}
	siteName := getSiteName(&doc, url)
	faviconURL := getFavicon(&doc, url, siteName)
	img := downloadAndDecodeImage(faviconURL)
	createFolder()
	iconPath := saveIco(siteName, img)
	createShortcut(url, siteName, iconPath)
}

func initIconDir() {
	userHome, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get user home directory: %v", err)
	}

	iconDir = filepath.Join(userHome, "Documents", "bureaubladiconen")
	fmt.Printf("Icon directory set to: %s\n", iconDir)
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

func downloadAndDecodeImage(faviconURL string) image.Image {
	iconRes, err := http.Get(faviconURL)
	if err != nil || iconRes.StatusCode != 200 {
		log.Fatalf("Failed to download favicon: %v (status: %d)", err, iconRes.StatusCode)
	}
	defer iconRes.Body.Close()

	// Read the entire response body into memory
	imageData, err := io.ReadAll(iconRes.Body)
	if err != nil {
		log.Fatalf("Failed to read image data: %v", err)
	}

	// Create a reader from the image data
	imageReader := strings.NewReader(string(imageData))

	img, format, err := image.Decode(imageReader)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}

	fmt.Printf("Successfully decoded favicon image (format: %s).\n", format)
	return img
}

func createFolder() {
	if _, err := os.Stat(iconDir); os.IsNotExist(err) {
		fmt.Printf("Creating directory: %s\n", iconDir)
		if err := os.MkdirAll(iconDir, 0755); err != nil {
			log.Fatalf("Failed to create icon directory: %v", err)
		}
	}
}

func saveIco(siteName string, img image.Image) string {
	safeFileName := strings.Map(func(r rune) rune {
		if strings.ContainsRune(`\/:*?"<>|`, r) {
			return '_'
		}
		return r
	}, siteName)

	iconPath := filepath.Join(iconDir, safeFileName+".ico")
	file, err := os.Create(iconPath)
	if err != nil {
		log.Fatalf("Failed to create .ico file: %v", err)
	}
	defer file.Close()

	if err := ico.Encode(file, img); err != nil {
		log.Fatalf("Failed to encode to .ico format: %v", err)
	}
	fmt.Printf("Icon saved to: %s\n", iconPath)

	return iconPath
}

func createShortcut(inputURL, siteName, iconPath string) {
	if err := createDesktopShortcut(inputURL, siteName, iconPath); err != nil {
		log.Fatalf("Failed to create shortcut: %v", err)
	}

	fmt.Println("\nShortcut created successfully on your desktop!")
}
