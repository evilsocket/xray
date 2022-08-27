package xray

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
)

type CertSH struct {
}

func NewCertSH() *CertSH {
	return &CertSH{}
}

func (me *CertSH) GetSubDomains(c *Context) []string {
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s", c.Domain)
	unique := make(map[string]bool)

	if res, err := http.Get(url); err == nil {
		defer res.Body.Close()

		if raw_body, err := ioutil.ReadAll(res.Body); err == nil {
			data := string(raw_body)
			re := regexp.MustCompile(fmt.Sprintf(">([^<\\*%%]+)\\.%s<", regexp.QuoteMeta(c.Domain)))
			if match := re.FindAllString(data, -1); match != nil {
				for _, m := range match {
					m = strings.Trim(m, "><")
					sub := c.GetSubDomain(m)
					unique[sub] = true
				}
			}
		}
	}

	sub := make([]string, len(unique))
	for k := range unique {
		sub = append(sub, k)
	}

	return sub
}
