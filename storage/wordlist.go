package storage

import (
	"bufio"
	"log"
	"os"

	"github.com/evilsocket/xray/core"
)

type Wordlist struct {
	Path string
	out  chan string
}

func loadWordlist(fileName string) (w *Wordlist, err error) {
	var fp *os.File

	w = &Wordlist{
		Path: fileName,
	}

	if core.Exists(w.Path) {
		log.Printf("loading wordlist from '%s'", w.Path)

		if fp, err = os.Open(w.Path); err != nil {
			return
		}

		w.out = make(chan string)
		go func() {
			defer fp.Close()
			defer close(w.out)
			scanner := bufio.NewScanner(fp)
			scanner.Split(bufio.ScanLines)
			for scanner.Scan() {
				w.out <- scanner.Text()
			}
		}()
	}

	return
}
