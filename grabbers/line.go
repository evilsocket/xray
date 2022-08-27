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
			defer func() {
				if err = conn.Close(); err != nil {
					fmt.Printf("error closing connection: %v\n", err)
				}
			}()

			if msg, err := bufio.NewReader(conn).ReadString('\n'); err == nil {
				t.Banners[g.Name()] = strings.Trim(msg, "\r\n\t ")
			}
		}
	}
}
