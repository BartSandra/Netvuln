package scanner

import (
	"context"
	"fmt"
	"github.com/Ullaakut/nmap"
	"github.com/sirupsen/logrus"
	pb "netvuln/api/proto"
)

// Scan выполняет сканирование с помощью Nmap.
func Scan(ctx context.Context, targets []string, port int32, log *logrus.Logger) (*pb.CheckVulnResponse, error) {
	// Логируем начало сканирования
	log.Infof("Starting scan for targets: %v on port %d", targets, port)

	// Создаем новый сканер Nmap
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(targets...),
		nmap.WithPorts(fmt.Sprintf("%d", port)),
		nmap.WithScripts("vulners"),
	)
	if err != nil {
		// Если не удалось создать сканер, логируем ошибку и возвращаем
		log.Errorf("Failed to create scanner: %v. Ensure Nmap is installed and available in $PATH.", err)
		return nil, fmt.Errorf("failed to create scanner: %w", err)
	}

	// Выполняем сканирование
	result, warnings, err := scanner.Run()
	if err != nil {
		log.Errorf("Failed to run scan: %v", err)
		return nil, fmt.Errorf("failed to run scan: %w", err)
	}

	// Логируем предупреждения, если они есть
	if warnings != nil {
		log.Warnf("Warnings: %v", warnings)
	}

	// Формируем ответ
	response := &pb.CheckVulnResponse{}
	for _, host := range result.Hosts {
		if len(host.Ports) == 0 {
			continue
		}

		log.Infof("Scanning host: %s", host.Addresses[0].Addr)

		targetResult := &pb.TargetResult{Target: host.Addresses[0].Addr}
		for _, port := range host.Ports {
			service := &pb.Service{
				Name:    port.Service.Name,
				Version: port.Service.Version,
				TcpPort: int32(port.ID),
			}

			log.Infof("Found service: %s, Version: %s, Port: %d", service.Name, service.Version, service.TcpPort)

			// Извлекаем уязвимости для найденных сервисов
			for _, script := range port.Scripts {
				if script.ID == "vulners" {
					for _, table := range script.Tables {
						for _, elem := range table.Elements {
							service.Vulns = append(service.Vulns, &pb.Vulnerability{
								Identifier: elem.Key,
								CvssScore:  parseCvssScore(elem.Value),
							})
						}
					}
				}
			}

			// Добавляем сервис в результаты
			targetResult.Services = append(targetResult.Services, service)
		}

		// Добавляем результаты для текущего хоста
		response.Results = append(response.Results, targetResult)
	}

	log.Infof("Scan completed successfully for targets: %v", targets)
	return response, nil
}

// parseCvssScore парсит CVSS-балл из строки
func parseCvssScore(value string) float32 {
	var score float32
	_, err := fmt.Sscanf(value, "%f", &score)
	if err != nil {
		logrus.Errorf("failed to parse CVSS score from value: %s, error: %v", value, err)
		return 0
	}
	return score
}
