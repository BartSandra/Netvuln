package service

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	pb "netvuln/api/proto"
	"netvuln/internal/config"
	"netvuln/internal/scanner"
)

// NetVulnService реализует NetVulnServiceServer.
type NetVulnService struct {
	pb.UnimplementedNetVulnServiceServer
}

// NewNetVulnService создает новый сервис.
func NewNetVulnService() *NetVulnService {
	return &NetVulnService{}
}

// CheckVuln выполняет сканирование уязвимостей.
func (s *NetVulnService) CheckVuln(ctx context.Context, req *pb.CheckVulnRequest) (*pb.CheckVulnResponse, error) {
	// Загружаем конфигурацию
	cfg := config.Load()

	// Создаем логгер для сервиса
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	// Настройка уровня логирования из конфигурации
	level, err := logrus.ParseLevel(cfg.LogLevel)
	if err != nil {
		log.Errorf("invalid log level %s, using default INFO", cfg.LogLevel)
		level = logrus.InfoLevel
	}
	log.SetLevel(level)

	log.Infof("Received CheckVuln request: targets=%v, tcp_port=%d", req.Targets, req.TcpPort)

	// Проверка наличия целей для сканирования
	if len(req.Targets) == 0 {
		log.Errorf("no targets provided")
		return nil, fmt.Errorf("no targets provided")
	}

	log.Infof("Starting scan for targets: %v on port %d", req.Targets, req.TcpPort)

	// Передаем в функцию Scan теперь context и log
	results, err := scanner.Scan(ctx, req.Targets, req.TcpPort, log)
	if err != nil {
		log.Errorf("failed to scan: %v", err)
		return nil, fmt.Errorf("failed to scan: %w", err)
	}

	log.Infof("Scan completed successfully: %+v", results)
	return results, nil
}
