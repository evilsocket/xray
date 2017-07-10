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
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"golang.org/x/net/html"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

type Dialer func(network, addr string) (net.Conn, error)

const DisallowLimit = 4

func isTitleElement(n *html.Node) bool {
	return n.Type == html.ElementNode && n.Data == "title"
}

func traverseHTML(n *html.Node) (string, bool) {
	if n != nil && isTitleElement(n) && n.FirstChild != nil {
		return n.FirstChild.Data, true
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		result, ok := traverseHTML(c)
		if ok {
			return result, ok
		}
	}

	return "", false
}

type HTTPGrabber struct {
}

func (g *HTTPGrabber) Name() string {
	return "http"
}

func makeDialer(certs *[]*x509.Certificate, skipCAVerification bool) Dialer {
	return func(network, addr string) (net.Conn, error) {
		c, err := tls.Dial(network, addr, &tls.Config{InsecureSkipVerify: skipCAVerification})
		if err != nil {
			return c, err
		}
		connstate := c.ConnectionState()
		*certs = connstate.PeerCertificates
		return c, nil
	}
}

func first(a []string) string {
	if len(a) > 0 {
		return a[0]
	}
	return ""
}

func Subject2String(s pkix.Name) string {
	return fmt.Sprintf("C=%s/O=%s/OU=%s/L=%s/P=%s/CN=%s",
		first(s.Country),
		first(s.Organization),
		first(s.OrganizationalUnit),
		first(s.Locality),
		first(s.Province),
		s.CommonName)
}

func collectCertificates(certs []*x509.Certificate, t *Target) {
	if certs != nil && len(certs) > 0 {
		for i, cert := range certs {
			t.Banners[fmt.Sprintf("https:chain[%d]", i)] = Subject2String(cert.Subject)
		}
	}
}

func collectHeaders(resp *http.Response, t *Target) {
	for name, value := range resp.Header {
		if name == "Server" {
			t.Banners["http:server"] = strings.Trim(value[0], "\r\n\t ")
		} else if name == "X-Powered-By" {
			t.Banners["http:poweredby"] = strings.Trim(value[0], "\r\n\t ")
		} else if name == "Location" {
			t.Banners["http:redirect"] = strings.Trim(value[0], "\r\n\t")
		}
	}
}

func collectHTML(resp *http.Response, t *Target) {
	if doc, err := html.Parse(resp.Body); err == nil {
		if title, found := traverseHTML(doc); found {
			t.Banners["html:title"] = strings.Trim(title, "\r\n\t ")
		}
	}
}

func collectRobots(client *http.Client, url string, t *Target) {
	rob, err := client.Get(url + "robots.txt")
	if err == nil {
		defer rob.Body.Close()
		if rob.StatusCode == 200 {
			raw, err := ioutil.ReadAll(rob.Body)
			if err == nil {
				data := string(raw)
				bann := make([]string, 0)

				for _, line := range strings.Split(data, "\n") {
					if strings.Contains(line, "Disallow:") {
						tok := strings.Trim(strings.Split(line, "Disallow:")[1], "\r\n\t ")
						if tok != "" {
							bann = append(bann, tok)
						}
					}

					if len(bann) >= DisallowLimit {
						bann = append(bann, "...")
						break
					}
				}

				if len(bann) > 0 {
					t.Banners["http:disallow"] = strings.Join(bann, ", ")
				}
			}
		}
	}
}

func (g *HTTPGrabber) Grab(port int, t *Target) {
	if port != 80 && port != 8080 && port != 443 && port != 8433 {
		return
	}

	base := ""
	url := ""

	if len(t.Domains) > 0 {
		base = t.Domains[0]
	} else {
		base = t.Address
	}

	if port == 80 {
		url = "http://" + base + "/"
	} else if port == 443 || port == 8433 {
		url = "https://" + base + "/"
	}

	certificates := make([]*x509.Certificate, 0)
	client := &http.Client{
		Transport: &http.Transport{
			DialTLS: makeDialer(&certificates, true),
		},
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.Get(url)
	if err == nil {
		defer resp.Body.Close()

		collectCertificates(certificates, t)
		collectHeaders(resp, t)
		collectHTML(resp, t)
	}

	collectRobots(client, url, t)
}
