package main

import (
	"context"
	"fmt"

	"dagger/tests/internal/dagger"
)

type Tests struct{}

func (m *Tests) Run(ctx context.Context, socket *dagger.Socket) error {
	if err := m.TestDefault(ctx, socket); err != nil {
		return fmt.Errorf("default test failed: %w", err)
	}

	if err := m.TestKubeConfig(ctx, socket); err != nil {
		return fmt.Errorf("kubeconfig test failed: %w", err)
	}

	return nil
}

func (m *Tests) TestDefault(ctx context.Context, socket *dagger.Socket) error {
	cluster := dag.Kind(socket).Cluster(dagger.KindClusterOpts{
		Name: "test-default",
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

func (m *Tests) TestKubeConfig(ctx context.Context, socket *dagger.Socket) error {
	customConfig := `kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
networking:
  disableDefaultCNI: true
`

	cluster := dag.Kind(socket).Cluster(dagger.KindClusterOpts{
		Name:   "test-kubeconfig",
		Config: customConfig,
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
