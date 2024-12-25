package main

import (
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	pb "netvuln/api/proto"
	"strings"
	"time"
)

func main() {
	// Чтение параметров командной строки
	address := flag.String("address", "localhost:50051", "Address of the gRPC server")
	targets := flag.String("targets", "", "Comma-separated list of target IPs or hostnames")
	port := flag.Int("port", 80, "TCP port to scan")
	flag.Parse()

	// Проверка наличия целей для сканирования
	if *targets == "" {
		log.Fatalf("No targets specified. Use the -targets flag to provide a comma-separated list of targets.")
	}

	// Разделение списка целей на массив строк
	targetList := strings.Split(*targets, ",")
	if len(targetList) == 0 {
		log.Fatalf("Invalid targets format. Ensure the -targets flag contains a valid comma-separated list of targets.")
	}

	// Установка соединения с gRPC сервером
	conn, err := grpc.NewClient(*address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to the gRPC server at %s: %v", *address, err)
	}
	defer conn.Close()

	// Создание клиента для взаимодействия с gRPC сервисом
	client := pb.NewNetVulnServiceClient(conn)

	// Формирование запроса на основе входных данных
	req := &pb.CheckVulnRequest{
		Targets: targetList,
		TcpPort: int32(*port),
	}

	// Установка тайм-аута для запроса
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Отправка запроса к сервису и получение ответа
	resp, err := client.CheckVuln(ctx, req)
	if err != nil {
		log.Fatalf("Error calling CheckVuln: %v", err)
	}

	// Вывод результатов сканирования
	if len(resp.Results) == 0 {
		log.Println("No vulnerabilities found for the specified targets.")
	} else {
		log.Println("Scan results:")
		for _, result := range resp.Results {
			fmt.Printf("Target: %s\n", result.Target)
			for _, service := range result.Services {
				fmt.Printf("  Service: %s (Version: %s, Port: %d)\n", service.Name, service.Version, service.TcpPort)
				if len(service.Vulns) > 0 {
					fmt.Println("    Vulnerabilities:")
					for _, vuln := range service.Vulns {
						fmt.Printf("      ID: %s, CVSS Score: %.1f\n", vuln.Identifier, vuln.CvssScore)
					}
				} else {
					fmt.Println("    No vulnerabilities found for this service.")
				}
			}
		}
	}
}
