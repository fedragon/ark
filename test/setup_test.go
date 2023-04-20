package test

import (
	"context"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	fmt.Println("Setting up test environment...")

	ctx := context.Background()
	natPort := "5432/tcp"
	req := testcontainers.ContainerRequest{
		Image:        "postgres:15-alpine",
		ExposedPorts: []string{natPort},
		Env: map[string]string{
			"POSTGRES_USER":     "ark",
			"POSTGRES_PASSWORD": "ark",
			"POSTGRES_DATABASE": "ark",
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections"),
			wait.ForListeningPort(nat.Port(natPort)),
		),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			log.Fatalf("failed to terminate container: %s", err.Error())
		}
	}()

	ip, err := container.Host(ctx)
	if err != nil {
		log.Fatal(err)
	}
	port, err := container.MappedPort(ctx, nat.Port(natPort))
	if err != nil {
		log.Fatal(err)
	}

	_ = os.Setenv("POSTGRES_ADDRESS", fmt.Sprintf("%s:%d", ip, port.Int()))

	code := m.Run()

	fmt.Println("Tearing down test environment...")

	_ = os.RemoveAll("archive")

	os.Exit(code)
}
