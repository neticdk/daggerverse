// A Dagger Module for integrating with the KinD
package main

import (
	"context"
	"dagger/kind/internal/dagger"
	"fmt"
)

const defaultImage = "alpine/k8s:1.35.0"

func New(
	// Unix socket to connect to the external Docker Engine. Please carefully use this option it can expose your host to the container.
	//
	// +required
	socket *dagger.Socket,
	// +optional
	containerImage string,
	// +optional
	kindImage string,
) *Kind {
	if containerImage == "" {
		containerImage = defaultImage
	}
	return &Kind{DockerSocket: socket, ContainerImage: containerImage, KindImage: kindImage}
}

type Kind struct {
	// +private
	DockerSocket *dagger.Socket
	// +private
	ContainerImage string
	// +private
	KindImage string
}

// Container that contains the kind and k9s binaries
func (k *Kind) Container() *dagger.Container {
	// Define the shim script content
	shim := `#!/bin/bash
args=("$@")

if [[ "${args[0]}" == "exec" ]]; then
    # Always use docker for exec to ensure pipe reliability
    cmd="/usr/bin/docker"
    new_args=()

    for arg in "${args[@]}"; do
        if [[ "$arg" == "kubeadm" ]]; then
            new_args+=("/usr/bin/setsid" "-w" "kubeadm")
        else
            new_args+=("$arg")
        fi
    done

    exec "$cmd" "${new_args[@]}"
else
    # Passthrough non-exec commands (run, inspect, network) to podman
    # Kind checks provider availability/version using these
    exec /usr/bin/podman.real "$@"
fi
`

	return dag.Container().
		From(k.ContainerImage).
		WithoutEntrypoint().
		WithUser("root").
		WithWorkdir("/").
		WithExec([]string{"apk", "add", "--no-cache", "podman", "kind", "k9s", "bash", "docker", "util-linux"}).
		WithExec([]string{"mv", "/usr/bin/podman", "/usr/bin/podman.real"}).
		WithNewFile("/usr/bin/podman", shim, dagger.ContainerWithNewFileOpts{Permissions: 0755}).
		WithUnixSocket("/var/run/docker.sock", k.DockerSocket).
		WithEnvVariable("CONTAINER_HOST", "unix:///var/run/docker.sock").
		WithEnvVariable("DOCKER_HOST", "unix:///var/run/docker.sock")
}

// Returns a cluster object that can be used to interact with the kind cluster
func (k *Kind) Cluster(
	ctx context.Context,

	// Name of the cluster
	//
	// +optional
	// +default="kind"
	name string,

	// If true, the default CNI is not used. This is useful for running kind clusters with a different CNI.
	//
	// +optional
	// +default=false
	disableDefaultCni bool,

	// Number of worker nodes for the Kind cluster.
	//
	// +optional
	// +default=0
	workerNodes int,
) (*Cluster, error) {
	targetNetwork := "kind"

	if err := connectEngineToNetwork(ctx, k.DockerSocket, targetNetwork); err != nil {
		return nil, fmt.Errorf("failed to multi-home dagger engine: %w", err)
	}

	return &Cluster{
		Name:              name,
		Network:           targetNetwork,
		Kind:              k,
		KindImage:         k.KindImage,
		DisableDefaultCni: disableDefaultCni,
		WorkerNodes:       workerNodes,
	}, nil
}
