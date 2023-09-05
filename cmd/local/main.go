package main

import (
	"log"
	"sync"

	"github.com/erdongli/shadowsocks-go/internal/cfg"
	"github.com/erdongli/shadowsocks-go/internal/tcp"
)

func main() {
	tcp, err := tcp.NewLocal(cfg.LocalPort, cfg.RemoteHost, cfg.RemotePort, cfg.PSK, cfg.AEADConfig)
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
