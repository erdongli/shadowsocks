package main

import (
	"flag"
	"log"
	"sync"

	"github.com/erdongli/shadowsocks-go/internal/shadow"
	"github.com/erdongli/shadowsocks-go/internal/tcp"
)

var (
	p = flag.String("p", "", "port to listen on")
	k = flag.String("k", "", "access key")
)

func main() {
	flag.Parse()
	if *p == "" || *k == "" {
		log.Fatal("missing port/access key")
	}

	tcp, err := tcp.NewRemote(*p, *k, shadow.ChaCha20Poly1305)
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
