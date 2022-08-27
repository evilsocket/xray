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
