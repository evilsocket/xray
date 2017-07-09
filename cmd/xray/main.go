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
package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/bobesa/go-domain-util/domainutil"
	"github.com/evilsocket/xray"
	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/ns3777k/go-shodan/shodan"
)

const version = "1.0.0b"

type Result struct {
	hostname string
	addrs    []string
}

func DoRequest(sub string) interface{} {
	hostname := fmt.Sprintf("%s.%s", sub, *base)
	if addrs, err := net.LookupHost(hostname); err == nil {
		return Result{hostname: hostname, addrs: addrs}
	}

	return nil
}

func OnResult(res interface{}) {
	result, ok := res.(Result)
	if !ok {
		fmt.Printf("Error while converting result.\n")
		return
	}

	for _, address := range result.addrs {
		// IPv4 only for now ¯\_(ツ)_/¯
		if strings.Contains(address, ":") {
			continue
		}

		t := pool.Find(address)
		if t == nil {
			pool.Add(xray.NewTarget(address, result.hostname, shapi))
			pool.FlushSession(&bruter.Stats)
		} else {
			if t.AddDomain(result.hostname) == true {
				pool.FlushSession(&bruter.Stats)
			}
		}
	}
}

var (
	session *xray.Session
	pool   *xray.Pool
	shapi  *shodan.Client
	bruter *xray.Machine
	router *gin.Engine

	base       = flag.String("domain", "", "Base domain to start enumeration from.")
	wordlist   = flag.String("wordlist", "wordlists/default.lst", "Wordlist file to use for enumeration.")
	consumers  = flag.Int("consumers", 16, "Number of concurrent consumers to use for subdomain enumeration.")
	shodan_tok = flag.String("shodan-key", "", "Shodan API key.")
	address    = flag.String("address", "127.0.0.1", "IP address to bind the web ui server to.")
	sesfile    = flag.String("session", SessionDefaultFilename, "Session file name.")
	port       = flag.Int("port", 8080, "TCP port to bind the web ui server to.")
)

func main() {
	flag.Parse()

	fmt.Println( "____  ___" )
	fmt.Println( "\\   \\/  /" )
	fmt.Println( " \\     RAY v", version )
	fmt.Println( " /    by Simone 'evilsocket' Margaritelli" )
	fmt.Println( "/___/\\  \\" )
	fmt.Println( "      \\_/" )
	fmt.Println( "" )

	if *base = domainutil.Domain(*base); *base == "" {
		fmt.Println("Invalid or empty domain specified.")
		flag.Usage()
		os.Exit(1)
	} else if *shodan_tok == "" {
		fmt.Printf("! WARNING: No Shodan API token provided, XRAY won't be able to get per-ip information.\n")
	} 
	
	if *sesfile == SessionDefaultFilename || *sesfile == "" {
		*sesfile = GetSessionFileName(*base )
	}

	gin.SetMode(gin.ReleaseMode)

	session = xray.NewSession(*sesfile)
	pool = xray.NewPool(session)
	shapi = shodan.NewClient(nil, *shodan_tok)
	bruter = xray.NewMachine(*consumers, *wordlist, session, DoRequest, OnResult)
	router = gin.New()

	// Easy stuff, serve static assets and JSON "API"
	router.Use(static.Serve("/", NewBFS("static")))
	router.GET("/targets", func(c *gin.Context) {
		bruter.UpdateStats()
		c.JSON(200, gin.H{
			"domain":  *base,
			"stats":   bruter.Stats,
			"targets": session.Targets,
		})
	})

	// Let the user know where the session file is located.
	if !pool.WasRestored() {
		fmt.Printf( "@ Saving session to %s\n", *sesfile )
	} else {
		fmt.Printf( "@ Restoring DNS bruteforcing from %.2f%%\n", session.Stats.Progress )
	}

	// Start web server in its own go routine.
	go func() {
		fmt.Printf("@ Web UI running on http://%s:%d/\n\n", *address, *port)
		if err := router.Run(fmt.Sprintf("%s:%d", *address, *port)); err != nil {
			panic(err)
		}
	}()

	// Save session and print progress every 10s.
	go func() {
		ticker := time.NewTicker(time.Millisecond * 10000)
		for _ = range ticker.C {
			bruter.UpdateStats()
			pool.FlushSession(&bruter.Stats)

			if bruter.Stats.Progress < 100.0 {
				fmt.Printf("%.2f %% completed, %.2f req/s, %d unique targets found so far ...\n", bruter.Stats.Progress, bruter.Stats.Eps, len(session.Targets))
			}
		}
	}()

	// Start DNS bruteforcing.
	if err := bruter.Start(); err != nil {
		panic(err)
	}

	bruter.Wait()

	fmt.Println("\nAll tasks completed, press Ctrl-C to quit.")

	// Wait forever ...
	select {}
}
