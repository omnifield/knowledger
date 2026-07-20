// Package config — env-only конфиг knowledger (маленькая Config с явным парсом, канон go.md).
// Ноль флагов/файлов: KNOWLEDGER_PORT (порт HTTP), KNOWLEDGER_DB (путь к БД-волюму).
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config — весь рантайм-ввод сервиса. Явные поля, явные дефолты, явная валидация.
type Config struct {
	Port int    // KNOWLEDGER_PORT — слушаем 0.0.0.0:Port (G1: bind wildcard, не loopback)
	DB   string // KNOWLEDGER_DB — путь к SQLite-файлу (per-product волюм); ":memory:" для тестов
}

const (
	// defaultPort — внутренний listening-порт knowledger (дверь /api/knowledger -> knowledger:8040).
	// Реальное закрепление номера — зона registry/devopser; тут дефолт.
	defaultPort = 8040
	defaultDB   = "knowledger.db"
)

// Load читает конфиг из окружения. Пустой/кривой ввод -> дефолт (кроме явно битого порта).
func Load() (Config, error) {
	c := Config{Port: defaultPort, DB: defaultDB}

	if v := strings.TrimSpace(os.Getenv("KNOWLEDGER_PORT")); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil || p < 1 || p > 65535 {
			return Config{}, fmt.Errorf("KNOWLEDGER_PORT=%q невалиден (ожидается 1..65535)", v)
		}
		c.Port = p
	}
	if v := strings.TrimSpace(os.Getenv("KNOWLEDGER_DB")); v != "" {
		c.DB = v
	}
	return c, nil
}

// Addr — адрес прослушивания. Bind 0.0.0.0 (G1): сосед по docker-сети должен достучаться,
// loopback -> 502 через дверь.
func (c Config) Addr() string {
	return fmt.Sprintf("0.0.0.0:%d", c.Port)
}
