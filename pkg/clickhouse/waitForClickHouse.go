package clickhouse

//import (
//	"database/sql"
//	"fmt"
//	"log"
//	"time"
//)
//
//func WaitForClickHouse(dsn string, maxRetries int, delay time.Duration) (*sql.DB, error) {
//	var db *sql.DB
//	var err error
//	for i := 0; i < maxRetries; i++ {
//		db, err = Connect(dsn)
//		if err == nil {
//			return db, nil
//		}
//		log.Printf("Не удалось подключиться к ClickHouse (попытка %d/%d): %v", i+1, maxRetries, err)
//		time.Sleep(delay)
//	}
//	return nil, fmt.Errorf("не удалось подключиться к ClickHouse после %d попыток: %w", maxRetries, err)
//}
