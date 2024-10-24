package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sync"
)

func main() {
	// Define CLI arguments
	valuesFile := flag.String("v", "", "Values CSV file")
	contentFile := flag.String("c", "", "Content CSV file")
	m2domain := flag.String("d", "", "M2 domain")
	outputDir := flag.String("o", "", "Output directory")
	parallelDownloads := flag.Int("p", 10, "Number of parallel downloads")
	flag.Parse()

	if (*valuesFile == "" && *contentFile == "") || *m2domain == "" || *outputDir == "" {
		fmt.Println("Magento 2 - Products Media Downloaded (CSV has to be prepared ahead)")
		fmt.Println("Usage example: mediadump ((-v file.csv) or (-c file.csv)) -d m2domain.com -o osprey")
		os.Exit(1)
	}

	files := make([]string, 0)
	if *contentFile != "" {
		files = append(files, mustExtractFiles(*contentFile)...)
	}

	if *valuesFile != "" {
		files = append(files, mustReadFiles(*valuesFile)...)
	}

	sem := make(chan struct{}, *parallelDownloads)
	var wg sync.WaitGroup

	for _, file := range files {
		wg.Add(1)

		go func(f string) {
			defer wg.Done()

			sem <- struct{}{}
			defer func() { <-sem }()

			err := downloadFile(*m2domain, f, *outputDir)
			if err != nil {
				fmt.Println("Error downloading file:", f, err)
			} else {
				fmt.Println("Successfully downloaded:", f)
			}
		}(file)
	}

	wg.Wait()
}

func downloadFile(domain, uri, dir string) error {
	// Construct the full URL
	url := fmt.Sprintf("https://%s%s", domain, uri)

	// Create the full output directory path by combining dir and uri path
	fullDir := filepath.Join(dir, filepath.Dir(uri))

	// Ensure the directory structure exists
	err := os.MkdirAll(fullDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create directories %s: %v", fullDir, err)
	}

	// Create the full output uri path
	fullPath := filepath.Join(dir, uri)

	// Open the uri for writing
	outFile, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("could not create uri %s: %v", fullPath, err)
	}
	defer outFile.Close()

	// Download the uri
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download uri from %s: %v", url, err)
	}
	defer resp.Body.Close()

	// Check if the server responded with a 200 OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Copy the downloaded content to the uri
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save uri to %s: %v", fullPath, err)
	}

	fmt.Printf("Successfully downloaded %s to %s\n", uri, fullPath)
	return nil
}

func removeDuplicate[T comparable](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}

func mustReadFiles(filePath string) []string {
	files, err := readFiles(filePath)
	if err != nil {
		panic(err)
	}
	return files
}

func readFiles(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %v", filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read all the records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %v", filePath, err)
	}

	// Collect all file names from the CSV
	var files []string
	isFirst := true
	for _, record := range records {
		if isFirst {
			isFirst = false
			continue
		}

		value := record[0]
		if len(record) > 0 {
			files = append(files, value)
		}
	}

	return removeDuplicate(files), nil
}

func mustExtractFiles(filePath string) []string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read all the records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	var files []string
	isFirst := true
	for _, record := range records {
		if isFirst {
			isFirst = false
			continue
		}

		content := record[0]
		if content != "" {
			files = append(files, extractMediaLinks(content)...)
		}
	}

	return removeDuplicate(files)
}

func extractMediaLinks(content string) []string {
	// Define a regex pattern to match the {{media url=...}} syntax
	re := regexp.MustCompile(`{{media url=([-_a-zA-Z0-9\.\/]+)}}`)

	// Find all matches
	matches := re.FindAllStringSubmatch(content, -1)

	// Create a slice to hold the converted links
	var mediaLinks []string

	// Loop over the matches and convert the media URLs
	for _, match := range matches {
		if len(match) > 1 {
			// Prepend "/media/" to the extracted URL part
			mediaLinks = append(mediaLinks, "/media/"+match[1])
		}
	}

	return mediaLinks
}
