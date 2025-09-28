package database

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"
)

// ConnectionPoolConfig holds database connection pool configuration
type ConnectionPoolConfig struct {
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
	ConnMaxIdleTime time.Duration
}

// DefaultConnectionPoolConfig returns optimized default settings for parallel processing
func DefaultConnectionPoolConfig() ConnectionPoolConfig {
	return ConnectionPoolConfig{
		MaxOpenConns:    100,              // Increased from 25 to 100 for parallel processing
		MaxIdleConns:    50,               // Increased from 10 to 50 to maintain warm connections
		ConnMaxLifetime: 30 * time.Minute, // Increased from 5 to 30 minutes for stability
		ConnMaxIdleTime: 5 * time.Minute,  // Increased from 1 to 5 minutes to reduce reconnection overhead
	}
}

// LoadConnectionPoolConfig loads configuration from environment variables with defaults
func LoadConnectionPoolConfig() ConnectionPoolConfig {
	config := DefaultConnectionPoolConfig()

	if val := os.Getenv("DB_MAX_OPEN_CONNS"); val != "" {
		if maxOpen, err := strconv.Atoi(val); err == nil {
			config.MaxOpenConns = maxOpen
		}
	}

	if val := os.Getenv("DB_MAX_IDLE_CONNS"); val != "" {
		if maxIdle, err := strconv.Atoi(val); err == nil {
			config.MaxIdleConns = maxIdle
		}
	}

	if val := os.Getenv("DB_CONN_MAX_LIFETIME_MINUTES"); val != "" {
		if lifetime, err := strconv.Atoi(val); err == nil {
			config.ConnMaxLifetime = time.Duration(lifetime) * time.Minute
		}
	}

	if val := os.Getenv("DB_CONN_MAX_IDLE_TIME_MINUTES"); val != "" {
		if idleTime, err := strconv.Atoi(val); err == nil {
			config.ConnMaxIdleTime = time.Duration(idleTime) * time.Minute
		}
	}

	return config
}

// Open creates a new database connection with optimized pool settings
func Open() (DB *sql.DB, err error) {
	DB, err = sql.Open("pgx", "host=localhost port=5432 user=movies password=movies dbname=movies sslmode=disable")
	if err != nil {
		return
	}

	// Load and apply connection pool configuration
	config := LoadConnectionPoolConfig()
	DB.SetMaxOpenConns(config.MaxOpenConns)
	DB.SetMaxIdleConns(config.MaxIdleConns)
	DB.SetConnMaxLifetime(config.ConnMaxLifetime)
	DB.SetConnMaxIdleTime(config.ConnMaxIdleTime)

	// Test the connection
	if err := DB.Ping(); err != nil {
		return nil, err
	}
	_, err = DB.Exec("DELETE FROM movies;")
	if err != nil {
		return
	}

	// Log connection pool configuration
	// fmt.Printf("Database connection pool configured: MaxOpen=%d, MaxIdle=%d, MaxLifetime=%v, MaxIdleTime=%v\n",
	// 	config.MaxOpenConns, config.MaxIdleConns, config.ConnMaxLifetime, config.ConnMaxIdleTime)

	return
}

// GetPoolStats returns current connection pool statistics for monitoring
func GetPoolStats(db *sql.DB) map[string]interface{} {
	stats := db.Stats()
	return map[string]interface{}{
		"max_open_connections": stats.MaxOpenConnections,
		"open_connections":     stats.OpenConnections,
		"in_use":               stats.InUse,
		"idle":                 stats.Idle,
		"wait_count":           stats.WaitCount,
		"wait_duration":        stats.WaitDuration,
		"max_idle_closed":      stats.MaxIdleClosed,
		"max_idle_time_closed": stats.MaxIdleTimeClosed,
		"max_lifetime_closed":  stats.MaxLifetimeClosed,
	}
}

// LogPoolStats logs current connection pool statistics
func LogPoolStats(db *sql.DB) {
	stats := GetPoolStats(db)
	fmt.Printf("Connection Pool Stats: Open=%d/%d, InUse=%d, Idle=%d, WaitCount=%d, WaitDuration=%v\n",
		stats["open_connections"], stats["max_open_connections"],
		stats["in_use"], stats["idle"], stats["wait_count"], stats["wait_duration"])
}
