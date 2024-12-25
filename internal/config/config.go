package config

import (
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var (
	cfg  *Config
	once sync.Once
)

// Config содержит конфигурацию сервиса.
type Config struct {
	Address  string
	LogLevel string
}

// Load загружает конфигурацию из переменных окружения.
// Конфигурация загружается только один раз, даже если функция вызывается несколько раз.
func Load() *Config {
	once.Do(func() {
		cfg = &Config{
			Address:  getEnv("SERVICE_ADDRESS", ":50051"), // Адрес сервиса (по умолчанию :50051)
			LogLevel: getEnv("LOG_LEVEL", "INFO"),         // Уровень логирования (по умолчанию INFO)
		}

		// Пример проверки корректности адреса
		if cfg.Address == "" {
			logrus.Fatalf("SERVICE_ADDRESS is not set or invalid")
		}

		// Устанавливаем уровень логирования
		setLogLevel(cfg.LogLevel)
	})
	return cfg
}

// getEnv получает значение переменной окружения по ключу.
// Если переменная окружения не установлена, возвращает значение по умолчанию.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	logrus.Warnf("Environment variable %s not set, using fallback value: %s", key, fallback)
	return fallback
}

// setLogLevel устанавливает уровень логирования в зависимости от значения переменной окружения.
func setLogLevel(level string) {
	switch level {
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
	case "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		logrus.Fatalf("Invalid LOG_LEVEL: %s. Valid levels are: DEBUG, INFO, WARN, ERROR, FATAL", level)
	}
}
