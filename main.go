package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"wayback_crawl/javascript" 
)

type WaybackResponse [][]string

func main() {
	domains := []string{"kifpool.me", "kifpool.app"}

	// Create the file to store URLs
	file, err := os.Create("urls.txt")
	if err != nil {
		fmt.Println("Error creating file:", err)
		os.Exit(1)
	}
	defer file.Close()

	// Map to store unique URLs
	uniqueUrls := make(map[string]struct{})

	// Fetch URLs for each domain
	for _, domain := range domains {
		file.WriteString("Results for domain: " + domain + "\n")

		// Run external commands to gather URLs
		runCommand(file, domain, "gau", []string{"--threads", "5"}, uniqueUrls)
		runCommand(file, domain, "waybackurls", []string{}, uniqueUrls)

		// Fetch URLs from web archive (Wayback Machine)
		url := fmt.Sprintf("https://web.archive.org/cdx/search/cdx?url=%s/*&output=json&collapse=urlkey", domain)
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("Error fetching data for domain", domain, ":", err)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Println("Error: Non-OK HTTP status for domain", domain, ":", resp.Status)
			continue
		}

		// Parse JSON response
		var data WaybackResponse
		if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
			fmt.Println("Error decoding JSON for domain", domain, ":", err)
			continue
		}

		// Save unique URLs from the web archive
		for _, entry := range data {
			if len(entry) > 2 {
				url := entry[2]
				if _, exists := uniqueUrls[url]; !exists && !isExcluded(url) {
					uniqueUrls[url] = struct{}{}
					file.WriteString(url + "\n")
				}
			}
		}

		file.WriteString("\n")
	}

	// Print success message
	fmt.Println("URLs saved to urls.txt successfully!")

	// Now, use the javascript package to find and print .js URLs
	err = javascript.FindJSURLs("urls.txt")
	if err != nil {
		fmt.Println("Error finding .js URLs:", err)
	}
}

func runCommand(file *os.File, domain, command string, args []string, uniqueUrls map[string]struct{}) {
	cmdArgs := append([]string{"echo", domain, "|", command}, args...)
	cmd := exec.Command(cmdArgs[0], cmdArgs[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing command:", err)
		return
	}

	for _, line := range strings.Split(string(output), "\n") {
		if line != "" && !isExcluded(line) {
			if _, exists := uniqueUrls[line]; !exists {
				uniqueUrls[line] = struct{}{}
				file.WriteString(line + "\n")
			}
		}
	}
}

func isExcluded(url string) bool {
	excludedFileTypes := []string{
		".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg", ".webp", 
		".css", ".scss", ".less", 
		".woff", ".woff2", ".ttf", ".otf", 
		".pdf", ".doc", ".docx", ".ppt", ".pptx", ".xls", ".xlsx",   
		".mp3", ".wav", ".flac", ".ogg", 
		".mp4", ".avi", ".mov", ".mkv", ".webm", 
		".ico", ".ico.ico", 
		".psd", ".ai", ".eps", ".indd", 
		".txt", ".md", ".rtf", "JPEG",
	}
	for _, ext := range excludedFileTypes {
		if strings.HasSuffix(url, ext) {
			return true
		}
	}
	return false
}

