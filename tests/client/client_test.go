package tests

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	pb "netvuln/api/proto"
	"testing"
)

func TestClient(t *testing.T) {
	server := grpc.NewServer()
	pb.RegisterNetVulnServiceServer(server, &mockNetVulnService{})

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	errCh := make(chan error)

	// Запуск сервера
	go func() {
		if err := server.Serve(listener); err != nil {
			errCh <- fmt.Errorf("failed to serve: %v", err)
		}
	}()
	defer server.Stop()

	// Подключение клиента
	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewNetVulnServiceClient(conn)

	req := &pb.CheckVulnRequest{
		Targets: []string{"192.168.100.14"},
		TcpPort: 80,
	}

	// Отправка запроса
	res, err := client.CheckVuln(context.Background(), req)
	if err != nil {
		t.Fatalf("could not check vulnerabilities: %v", err)
	}

	assert.NotNil(t, res)
	assert.Equal(t, len(res.Results), 1)
	assert.Equal(t, res.Results[0].Target, "192.168.100.14")
	assert.Equal(t, len(res.Results[0].Services), 1)
	assert.Equal(t, res.Results[0].Services[0].TcpPort, int32(80))

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("server error: %v", err)
		}
	default:
		// Все ок
	}
}

type mockNetVulnService struct {
	pb.UnimplementedNetVulnServiceServer
}

func (m *mockNetVulnService) CheckVuln(ctx context.Context, req *pb.CheckVulnRequest) (*pb.CheckVulnResponse, error) {
	return &pb.CheckVulnResponse{
		Results: []*pb.TargetResult{
			{
				Target: "192.168.100.14",
				Services: []*pb.Service{
					{
						Name:    "http",
						TcpPort: 80,
					},
				},
			},
		},
	}, nil
}
