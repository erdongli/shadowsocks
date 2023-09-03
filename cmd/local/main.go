package main

import (
	"log"
	"os"

	"github.com/erdongli/shadowsocks-go/internal/cfg"
	"github.com/erdongli/shadowsocks-go/internal/local"
)

func main() {
	log.SetOutput(os.Stdout)

	l, err := local.New(cfg.LocalPort, cfg.RemoteHost, cfg.RemotePort)
	if err != nil {
		log.Fatal(err)
	}

	if err := l.Serve(); err != nil {
		log.Fatal(err)
	}
}
