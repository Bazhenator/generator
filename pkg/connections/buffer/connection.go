package buffer

import (
	"context"

	"github.com/Bazhenator/tools/src/logger"
	"github.com/Bazhenator/generator/pkg/connections"
	"google.golang.org/grpc"
)

type Connection struct {
	CallOptions []grpc.CallOption

	conn *grpc.ClientConn
	l    *logger.Logger
}

func NewConnection(ctx context.Context, l *logger.Logger, target string) (*Connection, error) {
	conn, err := grpc.DialContext(ctx, target, connections.GetCommonDialOptions()...)
	if err != nil {
		l.Error("failed dial", logger.NewErrorField(err))
		return nil, err
	}

	res := &Connection{
		CallOptions: connections.CommonCallOptions,
		
		conn:        conn,
		l:           l,
	}
	return res, nil
}

func (c *Connection) Close() {
	if err := c.conn.Close(); err != nil {
		c.l.Error("failed close connection", logger.NewErrorField(err))
	}
}