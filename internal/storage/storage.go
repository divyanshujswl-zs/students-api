package storage

import (
	"fmt"

	"github.com/divyanshujswl-zs/students-api/internal/config"
	"github.com/divyanshujswl-zs/students-api/internal/storage/mysql"
	"github.com/divyanshujswl-zs/students-api/internal/storage/sqlite"
)

type Storage interface {
	CreateStudent(name, email string, age int) (int64, error)
}

func New(cfg *config.Config) (Storage, error) {

	switch cfg.DB.Driver {

	case "sqlite":
		return sqlite.New(cfg)

	case "mysql":
		return mysql.New(cfg)

	default:
		return nil, fmt.Errorf("unknown db driver: %s", cfg.DB.Driver)
	}
}
