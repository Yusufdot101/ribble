package grpc

import (
	"github.com/Yusufdot101/ripple-proto/golang/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	userClient user.UserServiceClient
	conn       *grpc.ClientConn
}

func NewAdapter(url string) (*Adapter, error) {
	conn, err := grpc.NewClient(
		url,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, err
	}

	client := user.NewUserServiceClient(conn)
	return &Adapter{
		conn:       conn,
		userClient: client,
	}, nil
}

func (a *Adapter) Close() error {
	return a.conn.Close()
}
