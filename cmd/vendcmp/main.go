package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
)

func main() {
	// Define CLI arguments
	projectPath := flag.String("p", "", "Path to the customized project")
	samplePath := flag.String("s", "", "Path to the non-customized project")
	flag.Parse()

	// Validate paths
	if *projectPath == "" || *samplePath == "" {
		fmt.Println("Both project and sample paths must be provided")
		os.Exit(1)
	}

	// Read vendor folders
	customizedVendors, err := readVendorDir(*projectPath)
	if err != nil {
		fmt.Println("Error reading customized project vendor directory:", err)
		os.Exit(1)
	}

	sampleVendors, err := readVendorDir(*samplePath)
	if err != nil {
		fmt.Println("Error reading sample project vendor directory:", err)
		os.Exit(1)
	}

	// Find additional packages
	additionalPackages := findAdditionalPackages(customizedVendors, sampleVendors)

	// Sort and print additional packages
	sort.Strings(additionalPackages)
	fmt.Println("Additional packages in the customized project:")
	for _, pkg := range additionalPackages {
		fmt.Println(pkg)
	}
}

// readVendorDir reads the vendor directory and returns a map of packages
func readVendorDir(path string) (map[string]bool, error) {
	vendorDir := filepath.Join(path, "vendor")
	files, err := ioutil.ReadDir(vendorDir)
	if err != nil {
		return nil, err
	}

	vendors := make(map[string]bool)
	for _, f := range files {
		if f.IsDir() {
			vendorPackages, err := ioutil.ReadDir(filepath.Join(vendorDir, f.Name()))
			if err != nil {
				return nil, err
			}
			for _, vp := range vendorPackages {
				if vp.IsDir() {
					vendors[f.Name()+"/"+vp.Name()] = true
				}
			}
		}
	}
	return vendors, nil
}

// findAdditionalPackages finds packages in customizedVendors not in sampleVendors
func findAdditionalPackages(customizedVendors, sampleVendors map[string]bool) []string {
	var additional []string
	for vendor, _ := range customizedVendors {
		if _, exists := sampleVendors[vendor]; !exists {
			additional = append(additional, vendor)
		}
	}
	return additional
}
