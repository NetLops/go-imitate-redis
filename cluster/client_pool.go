package cluster

import (
	"context"
	"errors"
	pool "github.com/jolestar/go-commons-pool/v2"
	"github.com/NetLops/go-imitate-redis/resp/client"
)

// 建立连接池
type connectionFactory struct {
	Peer string
}

func (c *connectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	connClient, err := client.MakeClient(c.Peer)
	if err != nil {
		return nil, err
	}
	connClient.Start()
	return pool.NewPooledObject(connClient), nil
}

func (c *connectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	connClient, ok := object.Object.(*client.Client)
	if !ok {
		return errors.New("type mismatch")
	}
	connClient.Close()
	return nil
}

func (c *connectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	return true
}

func (c *connectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}

func (c *connectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	return nil
}
