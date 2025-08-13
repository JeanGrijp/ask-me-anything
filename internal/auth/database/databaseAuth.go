package databaseAuth

// import (
// 	"context"
// 	"fmt"
// 	"os"
// 	"time"

// 	"github.com/JeanGrijp/petsync/internal/logger"

// 	"github.com/jmoiron/sqlx"
// 	_ "github.com/lib/pq"
// )

// func Connect() (*sqlx.DB, error) {
// 	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
// 	defer cancel()

// 	// use as variáveis genéricas, não JWT_DB_*
// 	host := os.Getenv("DB_HOST")
// 	port := os.Getenv("DB_PORT")
// 	user := os.Getenv("DB_USER")
// 	password := os.Getenv("DB_PASSWORD")
// 	dbname := os.Getenv("DB_NAME")
// 	timezone := os.Getenv("TZ")

// 	dsn := fmt.Sprintf(
// 		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s search_path=auth,core,public",
// 		host, port, user, password, dbname, timezone,
// 	)

// 	logger.Default.Info(ctx, "Connecting to database...", "dsn", dsn)
// 	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
// 	if err != nil {
// 		logger.Default.Error(ctx, "Failed to connect to database", "error", err)
// 		return nil, err
// 	}
// 	logger.Default.Info(ctx, "Connected to database successfully", "dsn", dsn)
// 	return db, nil
// }
