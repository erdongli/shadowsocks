package main

import (
	"flag"
	"log"
	"sync"

	"github.com/erdongli/shadowsocks/internal/shadow"
	"github.com/erdongli/shadowsocks/internal/tcp"
)

var (
	p = flag.String("p", "", "port to listen on")
	r = flag.String("r", "", "remote address to connect to")
	k = flag.String("k", "", "access key")
)

func main() {
	flag.Parse()
	if *p == "" || *r == "" || *k == "" {
		log.Fatal("missing port/remote address/access key")
	}

	tcp, err := tcp.NewLocal(*p, *r, *k, shadow.ChaCha20Poly1305)
	if err != nil {
		log.Fatal(err)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		tcp.Serve()
	}()

	wg.Wait()
}
