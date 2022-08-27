package grabbers

import (
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"regexp"
	"strings"

	xray "github.com/evilsocket/xray"

	"golang.org/x/net/html"
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

func collectCertificates(certs []*x509.Certificate, t *xray.Target) {
	if certs != nil && len(certs) > 0 {
		ctx := xray.GetContext()

		for i, cert := range certs {
			t.Banners[fmt.Sprintf("https:chain[%d]", i)] = Subject2String(cert.Subject)

			// Check for domains
			if ctx != nil {
				// Search in common name.
				if sub := ctx.GetSubDomain(cert.Subject.CommonName); sub != "" {
					ctx.Bruter.AddInput(sub)
				}

				// Search in alternative names.
				for _, name := range cert.DNSNames {
					if sub := ctx.GetSubDomain(name); sub != "" {
						ctx.Bruter.AddInput(sub)
					}
				}
			}
		}
	}
}

func collectHeaders(resp *http.Response, t *xray.Target) {
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

func collectHTML(resp *http.Response, t *xray.Target) {
	if raw_body, err := ioutil.ReadAll(resp.Body); err == nil {
		data := string(raw_body)

		// check if this is an Amazon bucket ... FUCK XML PARSERS!
		if strings.Contains(data, "ListBucketResult") && strings.Contains(data, "<Name>") {
			re := regexp.MustCompile(".*<Name>([^<]+)</Name>.*")
			match := re.FindStringSubmatch(data)
			if len(match) > 0 {
				t.Banners["amazon:bucket"] = match[1]
			}
		}

		// parse HTML for the title tag
		reader := strings.NewReader(data)
		if doc, err := html.Parse(reader); err == nil {
			if title, found := traverseHTML(doc); found {
				t.Banners["html:title"] = strings.Trim(title, "\r\n\t ")
			}
		}
	}
}

func collectRobots(client *http.Client, url string, t *xray.Target) {
	rob, err := client.Get(url + "robots.txt")
	if err != nil {
		return
	}

	defer rob.Body.Close()

	if rob.StatusCode != 200 {
		return
	}

	raw, err := ioutil.ReadAll(rob.Body)
	if err != nil {
		return
	}

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

func (g *HTTPGrabber) Grab(port int, t *xray.Target) {
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

	if port == 80 || port == 8080 {
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
