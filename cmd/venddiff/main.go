package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type VendorChange struct {
	Module  string
	FromVer string
	ToVer   string
}

type ChangeList struct {
	Current *VendorChange
	List    []*VendorChange
}

func (l *ChangeList) NewModule(moduleName string) {
	vendChang := &VendorChange{Module: moduleName}
	l.Current = vendChang
	l.List = append(l.List, vendChang)
}

func (l *ChangeList) AddFrom(fromVersion string) {
	l.Current.FromVer = fromVersion
}

func (l *ChangeList) AddTo(toVersion string) {
	l.Current.ToVer = toVersion
}

func NewChangeList() *ChangeList {
	return &ChangeList{
		List: make([]*VendorChange, 0),
	}
}

func main() {
	if len(os.Args) != 3 || os.Args[1] != "-d" {
		fmt.Println("Usage: venddiff -d <path_to_diff_file>")
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
	chLst := NewChangeList()
	isComposerLock := false

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "diff --git a/composer.lock b/composer.lock") {
			isComposerLock = true
		}

		if strings.HasPrefix(line, "diff --git") && !strings.Contains(line, "composer.lock") {
			isComposerLock = false
		}

		if !isComposerLock {
			continue
		}

		moduleName := ExtractValue(line, `"name":\s*"([^"]+)"`)
		if moduleName != "" {
			chLst.NewModule(moduleName)
		}

		fromVersion := ExtractValue(line, `^-[\s]*"version":\s*"([^"]+)"`)
		if fromVersion != "" {
			chLst.AddFrom(fromVersion)
		}

		toVersion := ExtractValue(line, `^\+[\s]*"version":\s*"([^"]+)"`)
		if toVersion != "" {
			chLst.AddTo(toVersion)
		}
	}

	for _, chg := range chLst.List {
		if chg.FromVer == "" || chg.ToVer == "" {
			continue
		}

		fmt.Println(fmt.Sprintf("%s: %s -> %s", chg.Module, chg.FromVer, chg.ToVer))
	}
}

func ExtractValue(line string, pattern string) string {
	re := regexp.MustCompile(pattern)
	matches := re.FindStringSubmatch(line)

	// Check if a match is found
	if len(matches) < 2 {
		return ""
	}

	return matches[1]
}
