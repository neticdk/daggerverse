package main

import (
	"context"
	"dagger/tests/internal/dagger"
	"fmt"
)

type Tests struct{}

func (m *Tests) Run(ctx context.Context, socket *dagger.Socket) error {
	cluster := dag.Kind(socket).Cluster(dagger.KindClusterOpts{
		Name:        "test",
		WorkerNodes: 2,
	})

	_, err := cluster.Create(ctx)
	if err != nil {
		return err
	}
	defer cluster.Delete(ctx)

	if exists, err := cluster.Exist(ctx); !exists || err != nil {
		if err != nil {
			return fmt.Errorf("checking if cluster exists: %w", err)
		}
		return fmt.Errorf("cluster does not exist")
	}

	return nil
}
