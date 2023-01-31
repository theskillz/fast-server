package stats

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/golang-migrate/migrate/v4"
	ch "github.com/golang-migrate/migrate/v4/database/clickhouse"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/theskillz/fast-server/internal/config"
)

type Stats interface {
	Run(ctx context.Context) error
	Add(useragent, ip string) error
	GetStatsForDay(day time.Time) ([]*Stat, error)
}

type service struct {
	chAddress  string
	chDatabase string
	chUsername string
	chPassword string
	conn       *sql.DB
}

type Stat struct {
	Useragent string
	IPAddress string
	Count     uint64
}

func NewStats(cfg *config.Config) Stats {
	return &service{
		chAddress:  cfg.Clickhouse.Address,
		chDatabase: cfg.Clickhouse.Database,
		chUsername: cfg.Clickhouse.Username,
		chPassword: cfg.Clickhouse.Password,
	}
}

func (s *service) Run(ctx context.Context) error {
	if err := s.openDB(); err != nil {
		return err
	}
	return nil
}

func (s *service) Add(ip, useragent string) error {
	_, err := s.conn.Exec("INSERT INTO stats(timestamp, useragent, ip_address) VALUES (now(), ?, ?)", useragent, ip)
	if err != nil {
		return err
	}
	return nil
}

func (s *service) GetStatsForDay(day time.Time) ([]*Stat, error) {
	result := make([]*Stat, 0)
	query := "SELECT useragent, ip_address, count(*) as count FROM stats WHERE toDate(timestamp) = ? GROUP BY toDate(timestamp), useragent, ip_address"
	rows, err := s.conn.Query(query, day.Format("2006-01-02"))
	if err != nil {
		return result, err
	}
	for rows.Next() {
		s := &Stat{}
		if err := rows.Scan(&s.Useragent, &s.IPAddress, &s.Count); err != nil {
			return result, err
		}
		result = append(result, s)
	}
	return result, nil
}

func (s *service) openDB() error {
	conn := clickhouse.OpenDB(&clickhouse.Options{
		Addr: []string{s.chAddress},
		Auth: clickhouse.Auth{
			Database: s.chDatabase,
			Username: s.chUsername,
			Password: s.chPassword,
		},
		Debug: false,
		Debugf: func(format string, v ...interface{}) {
			fmt.Printf(format, v)
		},
		Settings: clickhouse.Settings{
			"max_execution_time": 60,
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		DialTimeout:          time.Duration(1) * time.Second,
		ConnOpenStrategy:     clickhouse.ConnOpenInOrder,
		BlockBufferSize:      10,
		MaxCompressionBuffer: 10240,
	})
	conn.SetMaxIdleConns(20)
	conn.SetMaxOpenConns(20)
	conn.SetConnMaxLifetime(time.Duration(10) * time.Minute)
	s.conn = conn
	driver, err := ch.WithInstance(conn, &ch.Config{
		DatabaseName:    s.chDatabase,
		MigrationsTable: "migration",
	})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance("file://migrations", s.chDatabase, driver)
	if err != nil {
		return err
	}
	m.Up()
	return nil
}
