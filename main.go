package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/miekg/dns"
)

type Item struct {
	Ip     string
	Q      dns.Question
	Reject bool
	Answer []dns.RR
	Date   string
}

type CacheItem struct {
	Answer []dns.RR
	Time   time.Time
}

var BuildTime string

var queue chan Item

var cache map[dns.Question]CacheItem

var cli = new(dns.Client)

func parseQuery(m *dns.Msg, w dns.ResponseWriter) {

	for _, q := range m.Question {
		switch q.Qtype {
		case dns.TypeA:
			rej := !Validate(q.Name)

			clientip := strings.Split(w.RemoteAddr().String(), ":")[0]

			if !rej {
				m1 := new(dns.Msg)
				m1.Id = dns.Id()
				m1.RecursionDesired = true
				m1.Question = make([]dns.Question, 1)
				m1.Question = m.Question //dns.Question{"miek.nl.", dns.TypeMX, dns.ClassINET}

				if ans, ok := cache[q]; ok && time.Since(ans.Time).Seconds() < float64(ans.Answer[0].Header().Ttl) {
					m.Answer = ans.Answer
					i := Item{clientip, q, rej, nil, time.Now().Format(time.RFC3339)}
					queue <- i
					Log(i)
					// log.Println("found in cache")
				} else {
					relays := GetRelaydns()
					for _, dnsserver := range relays {
						in, _, err := cli.Exchange(m1, dnsserver)

						if err == nil {
							m.Answer = in.Answer
							i := Item{clientip, q, rej, in.Answer, time.Now().Format(time.RFC3339)}
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
				rip, err := GetRedirectIp()
				if err == nil {
					rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, rip))
					if err == nil {
						m.Answer = append(m.Answer, rr)
						Log(Item{clientip, q, rej, m.Answer, time.Now().Format(time.RFC3339)})
					}
				} else {
					Log(Item{clientip, q, rej, nil, time.Now().Format(time.RFC3339)})
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

		log.Printf("Q %s < %s : %t\n", x.Q.Name, x.Ip, x.Reject)

		Stat(x.Ip, x.Q.Name)

		if x.Answer != nil {
			log.Printf("ttl: %d \n", x.Answer[0].Header().Ttl)
			cache[x.Q] = CacheItem{x.Answer, time.Now()}
		}
	}
}

func main() {
	Init()

	cache = make(map[dns.Question]CacheItem)

	queue = make(chan Item, 100)

	go queueworker()

	// attach request handler func
	dns.HandleFunc(".", handleDnsRequest)

	// start server

	server := &dns.Server{Addr: config.DNSListen, Net: "udp"}

	if config.HttpListen != "" {
		log.Printf("Starting HTTP at port %v\n", config.HttpListen)
		go httpserver()
	}
	log.Printf("Starting DNS at port %v\n", config.DNSListen)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start DNS server: %s\n ", err.Error())
	}
}