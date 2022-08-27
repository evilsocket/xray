package xray

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var SessionDefaultFilename = "<domain-name>-xray-session.json"

type Session struct {
	filename string
	Stats    *Statistics
	Targets  map[string]*Target
}

func GetSessionFileName(domain string) string {
	return fmt.Sprintf("%s-xray-session.json", domain)
}

func NewSession(filename string) *Session {
	s := &Session{
		filename: filename,
		Stats:    nil,
		Targets:  make(map[string]*Target),
	}

	if _, err := os.Stat(s.filename); !os.IsNotExist(err) {
		fmt.Printf("@ Restoring session from %s ...\n", s.filename)
		if data, e := ioutil.ReadFile(s.filename); e == nil {
			if e = json.Unmarshal(data, &s); e != nil {
				panic(e)
			}

			fmt.Printf("@ Loaded %d entries from session file.\n", len(s.Targets))
		} else {
			panic(e)
		}
	}

	return s
}

func (s *Session) Flush(stats *Statistics) {
	s.Stats = stats
	if data, err := json.Marshal(s); err == nil {
		if err = ioutil.WriteFile(s.filename, data, 0644); err != nil {
			panic(err)
		}
	} else {
		panic(err)
	}
}
