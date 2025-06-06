package mock

import (
	"log"
	"net"

	"github.com/armon/go-socks5"
)

func StartMockSocks5Server(authEnabled bool) (addr string, closeFunc func()) {
	conf := &socks5.Config{}

	if authEnabled {
		conf.Credentials = socks5.StaticCredentials{"user": "pass"}
	}

	server, err := socks5.New(conf)
	if err != nil {
		log.Fatal(err)
	}

	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		_ = server.Serve(l)
	}()

	return l.Addr().String(), func() {
		_ = l.Close()
	}
}
