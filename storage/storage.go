package storage

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/evilsocket/xray/core"
)

const (
	DefaultPerms = 0700
	DefaultPath  = "~/.xray/"

	domainsFile = "domains.lst"
)

type Storage struct {
	sync.RWMutex

	Path    string
	Keys    *Keys
	Domains *Wordlist
}

var (
	I = (*Storage)(nil)
)

func Open(path string, create bool) (s *Storage, err error) {
	if path, err = core.ExpandPath(path); err != nil {
		return
	}

	if !core.Exists(path) {
		if create {
			log.Printf("creating folder '%s'", path)
			if err = os.MkdirAll(path, DefaultPerms); err != nil {
				return nil, err
			}
		} else {
			return nil, fmt.Errorf("path '%s' does not exist", path)
		}
	}

	s = &Storage{
		Path: path,
	}

	if s.Keys, err = loadKeys(path); err != nil {
		return nil, err
	}

	if s.Domains, err = loadWordlist(filepath.Join(path, domainsFile)); err != nil {
		return nil, err
	}

	I = s

	return
}
