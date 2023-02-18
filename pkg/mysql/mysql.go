// Package mysql implements mysql db connection.
package mysql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MySQLConfig - represents MySQL service config.
type MySQLConfig struct {
	User     string
	Password string
	Host     string
	Database string
}

// MySQL - represents mysql service.
type MySQL struct {
	DB *gorm.DB
}

// Model provides base fields for database models (like gorm.Model).
type Model struct {
	CreatedAt time.Time      `json:"createdAt" gorm:"index"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index" swaggerignore:"true"`
}

// Option - represents PostgreSQL service option.
type Option func(*MySQL)

// SetMaxIdleConns - configures max idle connections.
func SetMaxIdleConns(idleConns int) Option {
	return func(c *MySQL) {
		db, _ := c.DB.DB()
		db.SetMaxIdleConns(idleConns)
	}
}

// SetMaxOpenConns - consfigures max open connections.
func SetMaxOpenConns(openConns int) Option {
	return func(c *MySQL) {
		db, _ := c.DB.DB()
		db.SetMaxOpenConns(openConns)
	}
}

// SetConnMaxLifetime - configures max connection lifetime.
func SetConnMaxLifetime(maxLifetime time.Duration) Option {
	return func(c *MySQL) {
		db, _ := c.DB.DB()
		db.SetConnMaxLifetime(maxLifetime)
	}
}

// New - creates new instance of MySQL service.
func New(cfg MySQLConfig, opts ...Option) (*MySQL, error) {
	// create instance of mysql
	sql := &MySQL{}

	// apply custom options
	for _, opt := range opts {
		opt(sql)
	}

	// establish mysql connection
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.User, cfg.Password, cfg.Host, cfg.Database)
	sql.DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
	})

	if err != nil {
		return nil, fmt.Errorf("mysql - New - gorm.Open: %w", err)
	}

	return sql, nil
}

// Close - closes mysql service database connection.
func (p *MySQL) Close() {
	if p.DB != nil {
		db, _ := p.DB.DB()
		db.Close()
	}
}
