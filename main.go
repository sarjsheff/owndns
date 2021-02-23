package main

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/miekg/dns"
)

// Item - DNS request item.
type Item struct {
	IP     string
	Q      dns.Question
	Reject bool
	Answer []dns.RR
	Date   string
	Cache  bool
}

// CacheItem -
type CacheItem struct {
	Answer []dns.RR
	Time   time.Time
}

// BuildTime - Building time.
var BuildTime string

var queue chan Item

var synccache sync.Map

var cli = new(dns.Client)

func parseQuery(m *dns.Msg, w dns.ResponseWriter) {

	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:

			clientip := strings.Split(w.RemoteAddr().String(), ":")[0]

			rej := !Validate(q.Name, clientip)

			if !rej {
				m1 := new(dns.Msg)
				m1.Id = dns.Id()
				m1.RecursionDesired = true
				m1.Question = make([]dns.Question, 1)
				m1.Question = m.Question

				if ans, ok := synccache.Load(q); ok && time.Since(ans.(CacheItem).Time).Seconds() < float64(ans.(CacheItem).Answer[0].Header().Ttl) {
					m.Answer = ans.(CacheItem).Answer
					i := Item{clientip, q, rej, ans.(CacheItem).Answer, time.Now().Format(time.RFC3339), true}
					queue <- i
					Log(i)
					// log.Println("found in cache")
				} else {
					relays := sconfig.GetRelaydns()
					for _, dnsserver := range relays {
						in, _, err := cli.Exchange(m1, dnsserver.(string))

						if err == nil {
							m.Answer = in.Answer
							i := Item{clientip, q, rej, in.Answer, time.Now().Format(time.RFC3339), false}
							queue <- i
							Log(i)
							break
						} else {
							log.Println(dnsserver, err)
							//Log(Item{clientip, q, rej, nil})
						}
					}
				}
			} else {
				rip, err := GetRedirectIP()
				if err == nil {
					rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, rip))
					if err == nil {
						m.Answer = append(m.Answer, rr)
						Log(Item{clientip, q, rej, m.Answer, time.Now().Format(time.RFC3339), false})
					}
				} else {
					Log(Item{clientip, q, rej, nil, time.Now().Format(time.RFC3339), false})
				}
			}

		}
	}
}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m, w)
	}

	w.WriteMsg(m)
}

func queueworker() {
	for {
		x := <-queue

		log.Printf("Q %s < %s : %t\n", x.Q.Name, x.IP, x.Reject)

		Stat(x.IP, x.Q.Name)

		if x.Answer != nil {
			log.Printf("ttl: %d \n", x.Answer[0].Header().Ttl)
			tmp := CacheItem{x.Answer, time.Now()}
			synccache.Store(x.Q, tmp)
		}
	}
}

func main() {
	Init()

	queue = make(chan Item, 100)

	go queueworker()

	// attach request handler func
	dns.HandleFunc(".", handleDnsRequest)

	// start server

	server := &dns.Server{Addr: sconfig.Get("DNSListen").(string), Net: "udp"}

	if sconfig.Get("HTTPListen").(string) != "" {
		log.Printf("Starting HTTP at port %v\n", sconfig.Get("HTTPListen"))
		go httpserver()
	}
	log.Printf("Starting DNS at port %v\n", sconfig.Get("DNSListen"))
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start DNS server: %s\n ", err.Error())
	}
}
