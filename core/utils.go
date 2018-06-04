package core

import (
	"os"
	"os/user"
	"path/filepath"
	"sort"
	"strings"
)

const (
	defaultTrimSet = "\r\n\t "
)

func Trim(s string) string {
	return strings.Trim(s, defaultTrimSet)
}

func TrimRight(s string) string {
	return strings.TrimRight(s, defaultTrimSet)
}

func UniqueInts(a []int, sorted bool) []int {
	tmp := make(map[int]bool)
	uniq := make([]int, 0)

	for _, n := range a {
		tmp[n] = true
	}

	for n := range tmp {
		uniq = append(uniq, n)
	}

	if sorted {
		sort.Ints(uniq)
	}

	return uniq
}

func SepSplit(sv string, sep string) []string {
	filtered := make([]string, 0)
	for _, part := range strings.Split(sv, sep) {
		part = Trim(part)
		if part != "" {
			filtered = append(filtered, part)
		}
	}
	return filtered

}

func CommaSplit(csv string) []string {
	return SepSplit(csv, ",")
}

func Exists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func ExpandPath(path string) (string, error) {
	// Check if path is empty
	if path != "" {
		if strings.HasPrefix(path, "~") {
			usr, err := user.Current()
			if err != nil {
				return "", err
			} else {
				// Replace only the first occurrence of ~
				path = strings.Replace(path, "~", usr.HomeDir, 1)
			}
		}
		return filepath.Abs(path)
	}
	return "", nil
}
