package valvedproxy

import (
	"context"
	"log"
	"net"

	"github.com/filariow/gardenia/pkg/valvedprotos"
	"google.golang.org/grpc"
)

type Proxy interface {
	Serve(net.Listener) error
}

type proxy struct {
	valvedprotos.UnimplementedValvedSvcServer
	cli    valvedprotos.ValvedSvcClient
	server grpc.Server
}

func New(connectionString string) (Proxy, error) {
	// setup client
	conn, err := grpc.Dial(connectionString, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	cli := valvedprotos.NewValvedSvcClient(conn)

	// setup server
	s := grpc.NewServer()
	ps := proxy{cli: cli}
	valvedprotos.RegisterValvedSvcServer(s, &ps)

	return &proxy{cli: cli}, nil
}

// Open the valve
func (s *proxy) Open(ctx context.Context, req *valvedprotos.OpenValveRequest) (*valvedprotos.OpenValveReply, error) {
	log.Printf("Valve open request received")
	return s.cli.Open(ctx, req)
}

// Close the valve
func (s *proxy) Close(ctx context.Context, req *valvedprotos.CloseValveRequest) (*valvedprotos.CloseValveReply, error) {
	log.Printf("Valve open request received")
	return s.cli.Close(ctx, req)
}

// Returns the status of the valve
func (s *proxy) Status(ctx context.Context, req *valvedprotos.StatusValveRequest) (*valvedprotos.StatusValveReply, error) {
	log.Printf("Valve status request received")
	return s.cli.Status(ctx, req)
}

func (s *proxy) Serve(lis net.Listener) error {
	return s.server.Serve(lis)
}
