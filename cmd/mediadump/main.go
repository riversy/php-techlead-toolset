package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
)

func main() {
	// Define CLI arguments
	inputFile := flag.String("i", "", "Input CSV file")
	m2domain := flag.String("d", "", "M2 domain")
	outputDir := flag.String("o", "", "Output directory")
	parallelDownloads := flag.Int("p", 10, "Number of parallel downloads")
	flag.Parse()

	if *inputFile == "" || *m2domain == "" || *outputDir == "" {
		fmt.Println("Magento 2 - Products Media Downloaded (CSV has to be prepared ahead)")
		fmt.Println("Usage example: mediadump -i file.csv -d m2domain.com -o osprey")
		os.Exit(1)
	}

	files, err := readFiles(inputFile)
	if err != nil {
		fmt.Println("Error reading input file:", err)
		os.Exit(1)
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

func downloadFile(domain, file, dir string) error {
	// Construct the full URL
	url := fmt.Sprintf("https://%s/media/catalog/product%s", domain, file)

	// Create the full output directory path by combining dir and file path
	fullDir := filepath.Join(dir, filepath.Dir(file))

	// Ensure the directory structure exists
	err := os.MkdirAll(fullDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("could not create directories %s: %v", fullDir, err)
	}

	// Create the full output file path
	fullPath := filepath.Join(dir, file)

	// Open the file for writing
	outFile, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("could not create file %s: %v", fullPath, err)
	}
	defer outFile.Close()

	// Download the file
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download file from %s: %v", url, err)
	}
	defer resp.Body.Close()

	// Check if the server responded with a 200 OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Copy the downloaded content to the file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		return fmt.Errorf("failed to save file to %s: %v", fullPath, err)
	}

	fmt.Printf("Successfully downloaded %s to %s\n", file, fullPath)
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

func readFiles(filePath *string) ([]string, error) {
	file, err := os.Open(*filePath)
	if err != nil {
		return nil, fmt.Errorf("could not open file %s: %v", *filePath, err)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	// Read all the records from the CSV file
	records, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("could not read file %s: %v", *filePath, err)
	}

	// Collect all file names from the CSV
	var files []string
	for _, record := range records {
		value := record[0]
		if value == "value" {
			continue
		}

		if len(record) > 0 {
			files = append(files, value)
		}
	}

	return removeDuplicate(files), nil
}
