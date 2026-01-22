package mysql

import (
	"database/sql"
	"fmt"

	"github.com/divyanshujswl-zs/students-api/internal/config"
	_ "github.com/go-sql-driver/mysql"
)

type MySQL struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*MySQL, error) {

	// connect WITHOUT DB first
	rootDSN := fmt.Sprintf(
		"%s:%s@tcp(%s)/",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
	)

	db, err := sql.Open("mysql", rootDSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// create database
	_, err = db.Exec(fmt.Sprintf(
		"CREATE DATABASE IF NOT EXISTS %s",
		cfg.DB.Name,
	))
	if err != nil {
		return nil, err
	}

	db.Close()

	// reconnect WITH DB
	dbDSN := fmt.Sprintf(
		"%s:%s@tcp(%s)/%s?parseTime=true",
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Name,
	)

	db, err = sql.Open("mysql", dbDSN)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// 4️⃣ create table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
		id INT AUTO_INCREMENT PRIMARY KEY,
		name VARCHAR(255),
		email VARCHAR(255),
		age INT
	)`)

	if err != nil {
		return nil, err
	}

	return &MySQL{Db: db}, nil
}

func (m *MySQL) CreateStudent(name, email string, age int) (int64, error) {
	result, err := m.Db.Exec(
		"INSERT INTO students (name, email, age) VALUES (?, ?, ?)",
		name, email, age,
	)
	if err != nil {
		return 0, err
	}

	return result.LastInsertId()
}
