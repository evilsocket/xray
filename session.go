/*
 * Copyleft 2017, Simone Margaritelli <evilsocket at protonmail dot com>
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *   * Redistributions of source code must retain the above copyright notice,
 *     this list of conditions and the following disclaimer.
 *   * Redistributions in binary form must reproduce the above copyright
 *     notice, this list of conditions and the following disclaimer in the
 *     documentation and/or other materials provided with the distribution.
 *   * Neither the name of ARM Inject nor the names of its contributors may be used
 *     to endorse or promote products derived from this software without
 *     specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 */
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
