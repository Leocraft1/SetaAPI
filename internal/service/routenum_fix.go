package service

import (
	"setaapi/internal/repository"
	"sort"
	"unicode"
)

func GetRoutenums() []string {
	return sortLines(repository.GetRoutesDistinct())
}

func sortLines(lines []string) []string {
	sort.SliceStable(lines, func(i, j int) bool {
		numI := extractLineNumber(lines[i])
		numJ := extractLineNumber(lines[j])
		if numI != numJ {
			return numI < numJ
		}
		return lines[i] < lines[j]
	})

	numeric := make([]string, 0, len(lines))
	letters := make([]string, 0)

	for _, b := range lines {
		if len(b) > 0 && unicode.IsLetter(rune(b[0])) {
			letters = append(letters, b)
		} else {
			numeric = append(numeric, b)
		}
	}

	return append(numeric, letters...)
}