package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	pb "netvuln/api/proto"
	"netvuln/internal/config"
	"netvuln/internal/service"
)

func main() {
	// Настройка логирования
	log := logrus.New()
	log.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	// Загрузка конфигурации
	cfg := config.Load()

	// Уровень логирования из конфигурации
	switch cfg.LogLevel {
	case "DEBUG":
		log.SetLevel(logrus.DebugLevel)
	case "INFO":
		log.SetLevel(logrus.InfoLevel)
	case "ERROR":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	// Логирование конфигурации
	log.Infof("Loaded configuration: Address=%s, LogLevel=%s", cfg.Address, cfg.LogLevel)

	// Настройка gRPC-сервера
	server := grpc.NewServer()
	pb.RegisterNetVulnServiceServer(server, service.NewNetVulnService())

	// Настройка слушателя
	listener, err := net.Listen("tcp", cfg.Address)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Канал для обработки сигналов завершения
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Запуск сервера
	go func() {
		log.Infof("Starting server on %s", cfg.Address)
		if err := server.Serve(listener); err != nil {
			log.Fatalf("Server stopped with error: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-stop
	log.Info("Shutting down server...")

	// Корректное завершение работы сервера
	server.GracefulStop()
	listener.Close()
	log.Info("Server stopped gracefully")
}
