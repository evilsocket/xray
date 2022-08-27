package grabbers

import (
	"bufio"
	"fmt"
	"net"
	"regexp"

	xray "github.com/evilsocket/xray"
)

type MYSQLGrabber struct {
}

func (g *MYSQLGrabber) Name() string {
	return "mysql"
}

func (g *MYSQLGrabber) Grab(port int, t *xray.Target) {
	if port != 3306 {
		return
	}

	if conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", t.Address, port)); err == nil {
		defer func() {
			if err = conn.Close(); err != nil {
				fmt.Printf("error closing connection: %v\n", err)
			}
		}()

		buf := make([]byte, 1024)
		if read, err := bufio.NewReader(conn).Read(buf); err == nil && read > 0 {
			s := string(buf[0:read])
			re := regexp.MustCompile(".+\x0a([^\x00]+)\x00.+")
			match := re.FindStringSubmatch(s)
			if len(match) > 0 {
				t.Banners[g.Name()] = match[1]
			}
		}
	}
}
