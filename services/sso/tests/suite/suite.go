package suite

import (
	"context"
	"github.com/kordyd/remember_me-golang/protos/gens/go/sso"
	"github.com/kordyd/remember_me-golang/sso/internal/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

const grpcHost = "localhost"

type Suite struct {
	T          *testing.T
	AuthClient sso.AuthClient
	Cfg        *config.Config
}

const configPath = "../../config/config.yaml"

func New(t *testing.T) (*Suite, context.Context) {
	t.Helper()
	t.Parallel()
	cfg := config.MustLoadPath(configPath)
	ctx, cancelCtx := context.WithTimeout(context.Background(), cfg.GRPCConfig.Timeout)

	grpcAddress := net.JoinHostPort(grpcHost, strconv.Itoa(cfg.GRPCConfig.GRPCPort))

	conn, err := grpc.NewClient(grpcAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("could not create gRPC client: %v", err)
	}

	t.Cleanup(func() {
		t.Helper()
		cancelCtx()
	})

	return &Suite{
		T:          t,
		AuthClient: sso.NewAuthClient(conn),
		Cfg:        cfg,
	}, ctx
}
