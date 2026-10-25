package manager

import (
	"testing"
)

func TestRenameDuplicateFile(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		existing []string
		out      string
	}{
		{
			name:     "no collision",
			in:       "file.txt",
			existing: []string{},
			out:      "file.txt",
		},
		{
			name:     "name collision",
			in:       "file.txt",
			existing: []string{"file.txt"},
			out:      "file(1).txt",
		},
		{
			name:     "2nd name collision",
			in:       "file.txt",
			existing: []string{"file.txt", "file(1).txt"},
			out:      "file(2).txt",
		},
		{
			name:     "temporary name collision",
			in:       "file.txt",
			existing: []string{"file.txt.part"},
			out:      "file(1).txt",
		},
		{
			name:     "numbered temporary name collision",
			in:       "file.txt",
			existing: []string{"file.txt", "file(1).txt.part"},
			out:      "file(2).txt",
		},
		{
			name:     "multiple temporary name collisions",
			in:       "file.txt",
			existing: []string{"file.txt.part", "file(1).txt.part"},
			out:      "file(2).txt",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileSet := make(map[string]struct{})
			for _, name := range tt.existing {
				fileSet[name] = struct{}{}
			}
			fileName := RenameDuplicateFile(tt.in, fileSet)

			if fileName != tt.out {
				t.Errorf("got %s want %s", fileName, tt.out)
			}

			if _, ok := fileSet[fileName]; !ok {
				t.Errorf("fileName is missing  in set: %s", fileName)
			}
			if _, ok := fileSet[fileName+".part"]; !ok {
				t.Errorf("temporary filename %q is missing from set", fileName+".part")
			}
		})
	}
}

func TestRenameDuplicateFileReservesNames(t *testing.T) {
	fileSet := make(map[string]struct{})
	expected := []string{"file.txt", "file(1).txt", "file(2).txt"}

	for _, want := range expected {
		got := RenameDuplicateFile("file.txt", fileSet)
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		if _, ok := fileSet[got]; !ok {
			t.Errorf("filename %q is missing from set", got)
		}
		if _, ok := fileSet[got+".part"]; !ok {
			t.Errorf("temporary filename %q is missing from set", got+".part")
		}
	}
}
