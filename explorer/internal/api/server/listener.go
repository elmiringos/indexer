package server

import (
	"net"

	"github.com/soheilhy/cmux"
)

func SetupListeners(port string) (net.Listener, net.Listener, cmux.CMux, error) {
	l, err := net.Listen("tcp", port)
	if err != nil {
		return nil, nil, nil, err
	}

	m := cmux.New(l)

	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.HTTP1Fast())

	return grpcL, httpL, m, nil
}
