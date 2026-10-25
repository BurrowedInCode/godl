package manager

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func ExistingFilenames(directory string) (map[string]struct{}, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, fmt.Errorf("error reading directory: %w", err)
	}
	names := make(map[string]struct{}, len(entries))
	for _, entry := range entries {
		names[entry.Name()] = struct{}{}
	}
	return names, nil
}

func RenameDuplicateFile(file string, fileSet map[string]struct{}) string {
	fileNumber := 1
	tempFile := file + ".part"

	_, fileExists := fileSet[file]
	_, tempExists := fileSet[tempFile]

	if !fileExists && !tempExists {
		fileSet[file] = struct{}{}
		fileSet[tempFile] = struct{}{}
		return file
	}

	ext := filepath.Ext(file)
	fileName := strings.TrimSuffix(filepath.Base(file), ext)

	for {
		renamedDupeFile := fmt.Sprintf("%s(%d)%s", fileName, fileNumber, ext)
		renamedDupeTemp := renamedDupeFile + ".part"

		_, renamedDupeFileExists := fileSet[renamedDupeFile]
		_, renamedDupeTempExists := fileSet[renamedDupeTemp]

		if !renamedDupeFileExists && !renamedDupeTempExists {
			fileSet[renamedDupeFile] = struct{}{}
			fileSet[renamedDupeTemp] = struct{}{}
			return renamedDupeFile
		}
		fileNumber += 1
	}
}
