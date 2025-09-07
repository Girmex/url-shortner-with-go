package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
)

// URL struct represents each URL entry
type URL struct {
	Short string `json:"short"`
	Long  string `json:"long"`
}

// Global slice to store URLs
var urls []URL
const filename = "urls.json"

// Load URLs from JSON file
func loadURLs() {
	data, err := os.ReadFile(filename)
	if err != nil {
		urls = []URL{}
		return
	}
	err = json.Unmarshal(data, &urls)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		urls = []URL{}
	}
}

// Save URLs to JSON file
func saveURLs() {
	data, err := json.MarshalIndent(urls, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}
	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
	}
}

// Generate random 6-character short code
func generateShortCode() string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Shorten a URL and store it
func shorten(longURL string) string {
	short := generateShortCode()
	urls = append(urls, URL{Short: short, Long: longURL})
	saveURLs()
	return short
}

// Expand a short code to original URL
func expand(short string) string {
	for _, url := range urls {
		if url.Short == short {
			return url.Long
		}
	}
	return "URL not found!"
}

func main() {
	rand.Seed(time.Now().UnixNano())
	loadURLs()

	var rootCmd = &cobra.Command{
		Use:   "shortener",
		Short: "A simple URL shortener CLI",
	}

	var shortenCmd = &cobra.Command{
		Use:   "shorten [long-url]",
		Short: "Shorten a long URL",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			longURL := args[0]
			short := shorten(longURL)
			fmt.Println("Short URL code:", short)
		},
	}

	var expandCmd = &cobra.Command{
		Use:   "expand [short-code]",
		Short: "Expand a short URL code to the original URL",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			short := args[0]
			long := expand(short)
			fmt.Println("Original URL:", long)
		},
	}

	// Add subcommands
	rootCmd.AddCommand(shortenCmd)
	rootCmd.AddCommand(expandCmd)

	// Execute root command
	rootCmd.Execute()
}
