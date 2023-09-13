package main

import (
	"flag"
	"sync"

	"github.com/erdongli/shadowsocks/internal/log"
	"github.com/erdongli/shadowsocks/internal/shadow"
	"github.com/erdongli/shadowsocks/internal/tcp"
)

var (
	p = flag.String("p", "", "port to listen on")
	r = flag.String("r", "", "remote address to connect to")
	k = flag.String("k", "", "access key")
	l = flag.String("l", "info", "log level")
)

func main() {
	flag.Parse()

	log.SetLevel(*l)

	if *p == "" || *r == "" || *k == "" {
		log.Printf(log.Error, "missing port/remote address/access key")
		return
	}

	tcp, err := tcp.NewLocal(*p, *r, *k, shadow.ChaCha20Poly1305)
	if err != nil {
		log.Printf(log.Error, "failed to create TCP local: %v", err)
		return
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		tcp.Serve()
	}()

	wg.Wait()
}
