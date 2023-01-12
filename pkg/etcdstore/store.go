package etcdstore

import (
	"context"
	"fmt"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	dialTimeout    = 2 * time.Second
	requestTimeout = 10 * time.Second
)

type StoreConfig struct {
	URI string
}

type Store struct {
	config StoreConfig
	cli    *clientv3.Client
}

func (store *Store) Put(mac string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	lease, err := store.cli.Grant(ctx, 15)
	if err != nil {
		return err
	}

	key := fmt.Sprintf("/pxgclient/config/ignore/%s", mac)
	_, err = store.cli.Put(ctx, key, "x", clientv3.WithLease(lease.ID))
	if err != nil {
		return err
	}

	return nil
}

func NewStore(config StoreConfig) (*Store, error) {
	store := &Store{config: config}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"},
		DialTimeout: dialTimeout,
	})
	if err != nil {
		return nil, err
	}

	store.cli = cli
	return store, nil
}
