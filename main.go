package main

import (
	"fmt"
)

func main() {
	fmt.Println("URL Shortcut Creator")
	fmt.Println("What is the URL of the site you want to create a shortcut for?")
	var url string
	fmt.Scanln(&url)
	fmt.Println("> Creating shortcut for:", url)
}
