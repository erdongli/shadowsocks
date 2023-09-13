package tcp

import (
	"io"
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

		io.Copy(egress, ingress)
		egress.SetReadDeadline(time.Now().Add(timeout))
	}()

	io.Copy(ingress, egress)
	ingress.SetReadDeadline(time.Now().Add(timeout))

	wg.Wait()
}
