package changes

import (
	"fmt"
	"regexp"
)

type FileChangeType int

func (t *FileChangeType) String() string {
	switch *t {
	case ChangeTypeAdded:
		return "Added"
	case ChangeTypeModified:
		return "Modified"
	case ChangeTypeRemoved:
		return "Removed"
	case ChangeTypeMoved:
		return "Moved"
	default:
		return "Unknown"
	}
}

const (
	ChangeTypeUndefined = FileChangeType(iota)
	ChangeTypeAdded
	ChangeTypeModified
	ChangeTypeMoved
	ChangeTypeRemoved
)

type FileChangeBuilder struct {
	aPart      string
	bPart      string
	fileChange *FileChange
}

func NewFileChangeBuilder() *FileChangeBuilder {
	return &FileChangeBuilder{
		fileChange: &FileChange{},
	}
}

func (b *FileChangeBuilder) WithFilePath(filePath string) *FileChangeBuilder {
	b.fileChange.FilePath = filePath
	return b
}

func (b *FileChangeBuilder) WithAPart(aPart string) *FileChangeBuilder {
	b.aPart = aPart
	return b
}

func (b *FileChangeBuilder) WithBPart(bPart string) *FileChangeBuilder {
	b.bPart = bPart
	return b
}

func (b *FileChangeBuilder) defineChangeType() {
	changeType := ChangeTypeUndefined

	if b.aPart == b.bPart {
		changeType = ChangeTypeModified
	} else {
		changeType = ChangeTypeMoved
	}

	if b.aPart == "" && b.bPart != "" {
		changeType = ChangeTypeAdded
	}

	if b.bPart == "" && b.aPart != "" {
		changeType = ChangeTypeRemoved
	}

	b.fileChange.ChangeType = changeType
}

func (b *FileChangeBuilder) Build() *FileChange {
	b.defineChangeType()
	return b.fileChange
}

type FileChange struct {
	FilePath   string
	ChangeType FileChangeType
}

func (f *FileChange) CanPrint() bool {
	return f.ChangeType != ChangeTypeUndefined
}

func (f *FileChange) String() string {
	return fmt.Sprintf("%s - %s", f.FilePath, f.ChangeType.String())
}

type FileChanges struct {
	Changes []*FileChange
}

func (c *FileChanges) AddChange(fileChange *FileChange) {
	c.Changes = append(c.Changes, fileChange)
}

func NewFileChanges() *FileChanges {
	return &FileChanges{
		Changes: make([]*FileChange, 0),
	}
}

func DoesMatchPattern(line string, pattern string) bool {
	regex := regexp.MustCompile(pattern)
	return regex.MatchString(line)
}

func ExtractByPattern(line string, pattern string) string {
	regex := regexp.MustCompile(pattern)
	matches := regex.FindStringSubmatch(line)
	if len(matches) < 2 {
		return ""
	}
	return matches[1]
}
