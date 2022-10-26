package main

import (
	dbsql "database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type sql struct{}
type keyValue map[string]interface{}

func New() *sql {
	return &sql{}
}

func contains(array []string, element string) bool {
	for _, item := range array {
		if item == element {
			return true
		}
	}
	return false
}

func (*sql) Open(database string, connectionString string) *dbsql.DB {
	supportedDatabases := []string{"mysql", "postgres", "sqlite3"}
	if !contains(supportedDatabases, database) {
		log.Fatal("Database is not supported")
		return nil
	}

	db, err := dbsql.Open(database, connectionString)
	if err == nil {
		return db
	}

	log.Fatal(err)
	return nil
}

func (*sql) Query(db *dbsql.DB, query string) []keyValue {
	rows, _ := db.Query(query)
	cols, _ := rows.Columns()
	values := make([]interface{}, len(cols))
	valuePtrs := make([]interface{}, len(cols))
	result := make([]keyValue, 0)

	for rows.Next() {
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		err := rows.Scan(valuePtrs...)

		if err != nil {
			log.Fatal(err)
			return nil
		}

		data := make(keyValue, len(cols))
		for i, colName := range cols {
			data[colName] = *valuePtrs[i].(*interface{})
		}
		result = append(result, data)
	}

	rows.Close()
	return result
}
