package main

import (
	"log"
	"os"

	"github.com/erdongli/shadowsocks-go/internal/cfg"
	"github.com/erdongli/shadowsocks-go/internal/remote"
)

func main() {
	log.SetOutput(os.Stdout)

	r, err := remote.New(cfg.RemotePort)
	if err != nil {
		log.Fatal(err)
	}

	if err := r.Serve(); err != nil {
		log.Fatal(err)
	}
}
