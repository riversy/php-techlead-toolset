package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/riversy/php-techlead-toolset/pkg/changes"
)

var changeBuilder *changes.FileChangeBuilder

func main() {

	if len(os.Args) != 3 || os.Args[1] != "-d" {
		fmt.Println("Usage: filesdiff -d <path_to_diff_file>")
		os.Exit(1)
	}

	diffFilePath := os.Args[2]
	file, err := os.Open(diffFilePath)
	if err != nil {
		fmt.Printf("Error opening file: %s\n", err)
		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	fileChanges := changes.NewFileChanges()
	changeBuilder = nil

	for scanner.Scan() {
		line := scanner.Text()

		if isNewFileChange(line) {
			if changeBuilder != nil {
				fileChanges.AddChange(changeBuilder.Build())
			}

			filePath := extractFileChange(line)
			changeBuilder = changes.NewFileChangeBuilder().WithFilePath(filePath)
		}

		if isAPart(line) {
			aPart := extractAPart(line)
			changeBuilder.WithAPart(aPart)
		}

		if isBPart(line) {
			bPart := extractBPart(line)
			changeBuilder.WithBPart(bPart)
		}
	}

	if changeBuilder != nil {
		fileChanges.AddChange(changeBuilder.Build())
	}

	for _, change := range fileChanges.Changes {
		fmt.Println(change)
	}
}

func extractBPart(line string) string {
	return changes.ExtractByPattern(line, `^\+\+\+ b/(.*)`)
}

func isBPart(line string) bool {
	return changes.DoesMatchPattern(line, `^\+\+\+ b/.*`)
}

func extractAPart(line string) string {
	return changes.ExtractByPattern(line, `^--- a/(.*)`)
}

func isAPart(line string) bool {
	return changes.DoesMatchPattern(line, `^--- a/.*`)
}

func extractFileChange(line string) string {
	return changes.ExtractByPattern(line, `^diff --git a/(.*) b/.*$`)
}

func isNewFileChange(line string) bool {
	return changes.DoesMatchPattern(line, `^diff --git a/.* b/.*$`)
}
