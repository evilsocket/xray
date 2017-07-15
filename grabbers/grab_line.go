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
package grabbers

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	xray "github.com/evilsocket/xray"
)

type LineGrabber struct {
	name  string
	ports []int
}

func NewLineGrabber(name string, ports []int) *LineGrabber {
	return &LineGrabber{
		name:  name,
		ports: ports,
	}
}

func (g *LineGrabber) Name() string {
	return g.name
}

func (g *LineGrabber) CheckPort(port int) bool {
	for _, p := range g.ports {
		if p == port {
			return true
		}
	}
	return false
}

func (g *LineGrabber) Grab(port int, t *xray.Target) {
	if g.CheckPort(port) {
		if conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", t.Address, port)); err == nil {
			/*empijei: error not handled,
			suggestion to either handle it or throw it away properly, i.e.:
			defer func(){
				_ = conn.Close()
			}
			*/
			defer conn.Close()
			msg, _ := bufio.NewReader(conn).ReadString('\n')
			t.Banners[g.Name()] = strings.Trim(msg, "\r\n\t ")
		}
	}
}
