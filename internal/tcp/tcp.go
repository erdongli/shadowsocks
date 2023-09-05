package tcp

import (
	"io"
	"log"
	"net"
	"sync"
	"time"
)

const (
	network = "tcp"
	timeout = 5 * time.Second
)

func relay(egress, ingress net.Conn) {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()

		if _, err := io.Copy(egress, ingress); err != nil {
			log.Println(err)
		}

		egress.SetReadDeadline(time.Now().Add(timeout))
	}()

	io.Copy(ingress, egress)
	ingress.SetReadDeadline(time.Now().Add(timeout))

	wg.Wait()
}
