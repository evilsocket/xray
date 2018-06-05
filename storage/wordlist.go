package storage

import (
	"bufio"
	"log"
	"os"

	"github.com/evilsocket/xray/core"
)

type Wordlist []string

func loadWordlist(fileName string) (w Wordlist, err error) {
	var fp *os.File

	w = make(Wordlist, 0)
	tmp := make(map[string]bool)

	if core.Exists(fileName) {
		log.Printf("loading wordlist from '%s'", fileName)

		if fp, err = os.Open(fileName); err != nil {
			return
		}
		defer fp.Close()

		scanner := bufio.NewScanner(fp)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			word := core.Trim(scanner.Text())
			if found, _ := tmp[word]; !found {
				tmp[word] = true
				w = append(w, word)
			}
		}
	}

	return
}
