package main

import (
	"log"

	"github.com/erdongli/shadowsocks-go/internal/cfg"
	"github.com/erdongli/shadowsocks-go/internal/local"
)

func main() {
	l, err := local.New(cfg.Port)
	if err != nil {
		log.Fatal(err)
	}

	if err := l.Serve(); err != nil {
		log.Fatal(err)
	}
}
