package main

import (
	"context"
	"fmt"

	"dagger/tests/internal/dagger"
)

type Tests struct {
	cli *dagger.DockerCli
}

func (m *Tests) Run(ctx context.Context) error {
	cli := dag.Docker().Cli()
	m.cli = cli

	if err := m.TestCliVersionCmd(ctx); err != nil {
		return err
	}

	return nil
}

func (m *Tests) TestCliVersionCmd(ctx context.Context) error {
	if _, err := m.cli.Run(ctx, []string{"version"}, dagger.DockerCliRunOpts{}); err != nil {
		return fmt.Errorf("running CLI command 'version': %w", err)
	}
	return nil
}
