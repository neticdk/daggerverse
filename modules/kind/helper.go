package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"dagger/kind/internal/dagger"
)

// Helper Functions

func getContainerNetwork(ctx context.Context, socket *dagger.Socket, name string) (string, error) {
	return dag.Docker().
		Cli(dagger.DockerCliOpts{Socket: socket}).
		Run(ctx, []string{"container", "ls", "--filter", fmt.Sprintf("name=%s", name), "--format", "{{.Networks}}", "-n", "1"}, dagger.DockerCliRunOpts{InvalidateCache: true})
}

func exec(ctx context.Context, ctr *dagger.Container, kindNetwork string, args ...string) (*dagger.Container, error) {
	return ctr.
		WithEnvVariable("KIND_EXPERIMENTAL_DOCKER_NETWORK", kindNetwork).
		WithEnvVariable("CACHE_BUSTER", time.Now().Format(time.RFC3339Nano)).
		WithExec(args).
		Sync(ctx)
}

func getClusterIPAddress(ctx context.Context, socket *dagger.Socket, network, name string) (string, error) {
	return dag.Docker().
		Cli(dagger.DockerCliOpts{Socket: socket}).
		Run(ctx, []string{"inspect", fmt.Sprintf("%s-control-plane", name), "--format", fmt.Sprintf("{{.NetworkSettings.Networks.%s.IPAddress}}", network)}, dagger.DockerCliRunOpts{InvalidateCache: true})
}

// Connects the Dagger Engine container to a specific network
func connectEngineToNetwork(ctx context.Context, socket *dagger.Socket, networkName string) error {
	engineName, err := dag.Docker().
		Cli(dagger.DockerCliOpts{Socket: socket}).
		Run(ctx, []string{"container", "ls", "--filter", "name=^dagger-engine-.*", "--format", "{{.ID}}", "-n", "1"}, dagger.DockerCliRunOpts{InvalidateCache: true})

	if err != nil {
		return fmt.Errorf("failed to find dagger engine container: %w", err)
	}
	engineName = strings.TrimSpace(engineName)
	if engineName == "" {
		return fmt.Errorf("no running dagger engine found")
	}

	// Create the network if it doesn't exist (e.g. if docker and not podman...)
	_, _ = dag.Docker().
		Cli(dagger.DockerCliOpts{Socket: socket}).
		Run(ctx, []string{"network", "create", networkName})

	_, err = dag.Docker().
		Cli(dagger.DockerCliOpts{Socket: socket}).
		Run(ctx, []string{"network", "connect", networkName, engineName})

	if err != nil && !strings.Contains(err.Error(), "already exists") && !strings.Contains(err.Error(), "already connected") {
		fmt.Printf("Warning: attempted to connect engine to network %s but got: %v\n", networkName, err)
	}

	return nil
}
