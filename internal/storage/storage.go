package storage

import (
	"database/sql"
	"reflect"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
)

type Storage struct {
	db      *sql.DB
	dialect string
}

func New(db *sql.DB, dialect string) *Storage {
	return &Storage{db: db, dialect: dialect}
}

// It must take pointer to the structure.
func getColumnsForStruct(data interface{}) []interface{} {
	s := reflect.ValueOf(data).Elem()
	numCols := s.NumField()
	columns := make([]interface{}, numCols)

	for i := 0; i < numCols; i++ {
		field := s.Field(i)
		columns[i] = field.Addr().Interface()
	}

	return columns
}

//later for testing
// func getColumnsForStruct(structs []interface{}) []interface{} {
// 	result := []interface{}{}
// 	for _, data := range structs {
// 		s := reflect.ValueOf(data).Elem()
// 		numCols := s.NumField()
// 		columns := make([]interface{}, numCols)

// 		for i := 0; i < numCols; i++ {
// 			field := s.Field(i)
// 			columns[i] = field.Addr().Interface()
// 		}

// 		result = append(result, columns...)
// 	}

// 	return result
// }
