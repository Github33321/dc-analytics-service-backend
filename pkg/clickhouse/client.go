package clickhouse

//
//import (
//	"database/sql"
//	_ "github.com/ClickHouse/clickhouse-go"
//)
//
//func Connect(dsn string) (*sql.DB, error) {
//	db, err := sql.Open("clickhouse", dsn)
//	if err != nil {
//		return nil, err
//	}
//	if err := db.Ping(); err != nil {
//		return nil, err
//	}
//	return db, nil
//}
