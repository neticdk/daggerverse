package main

import (
	"context"

	"dagger/tests/internal/dagger"
)

type Tests struct{}

func (m *Tests) Run(ctx context.Context, socket *dagger.Socket) error {
	cluster := dag.Kind(socket).Cluster(dagger.KindClusterOpts{
		Name: "test",
	})

	_, err := cluster.Create(ctx)
	if err != nil {
		return err
	}
	// defer cluster.Delete(ctx)

	return nil
}
