package storage

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"

	"github.com/evilsocket/xray/core"
)

const (
	keysFilename = "keys.yml"
)

type Keys struct {
	sync.RWMutex
	Path string
	data map[string]string
}

func loadKeys(path string) (k *Keys, err error) {
	k = &Keys{
		Path: filepath.Join(path, keysFilename),
		data: make(map[string]string),
	}

	if core.Exists(k.Path) {
		log.Printf("loading keys from '%s'", k.Path)

		var raw []byte
		if raw, err = ioutil.ReadFile(k.Path); err != nil {
			return
		}

		err = yaml.Unmarshal(raw, &k.data)

		for k, _ := range k.data {
			log.Printf("loaded '%s' API key", k)
		}
	}

	return
}

func (k *Keys) Have(name string) bool {
	k.RLock()
	defer k.RUnlock()

	if v, found := k.data[name]; found {
		return core.Trim(v) != ""
	}
	return false
}
