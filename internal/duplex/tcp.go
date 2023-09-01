package duplex

import (
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	timeout = 60 * time.Second
)

func Relay(ingress, egress net.Conn) {
	log.Printf("relaying between %s and %s", ingress.LocalAddr(), egress.RemoteAddr())

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()

		io.Copy(egress, ingress)
	}()

	io.Copy(ingress, egress)

	wg.Wait()
}
