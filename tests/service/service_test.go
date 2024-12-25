package tests

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	pb "netvuln/api/proto"
	"netvuln/internal/service"
	"testing"
)

func TestCheckVuln(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.InfoLevel)

	server := grpc.NewServer()
	pb.RegisterNetVulnServiceServer(server, service.NewNetVulnService())

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}

	errCh := make(chan error)

	// Запускаем сервер
	go func() {
		if err := server.Serve(listener); err != nil {
			errCh <- fmt.Errorf("failed to serve: %v", err)
		}
	}()
	defer server.Stop()

	// Подключаем клиента
	conn, err := grpc.NewClient(":50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewNetVulnServiceClient(conn)

	// Формируем запрос
	req := &pb.CheckVulnRequest{
		Targets: []string{"8.8.8.8"}, // Используем публичный IP-адрес (например, DNS-сервер Google)
		TcpPort: 80,
	}

	// Отправляем запрос
	res, err := client.CheckVuln(context.Background(), req)
	if err != nil {
		t.Fatalf("could not check vulnerabilities: %v", err)
	}

	assert.NotNil(t, res, "Response should not be nil")
	assert.Greater(t, len(res.Results), 0, "Expected at least 1 result, but got none")

	if len(res.Results) > 0 {
		assert.Equal(t, "8.8.8.8", res.Results[0].Target, "Unexpected target")
		assert.Greater(t, len(res.Results[0].Services), 0, "Expected at least 1 service")
		if len(res.Results[0].Services) > 0 {
			assert.Equal(t, int32(80), res.Results[0].Services[0].TcpPort, "Unexpected port")
		}
	}

	select {
	case err := <-errCh:
		t.Fatalf("server error: %v", err)
	default:
		// Все ок
	}
}
