package storage

import (
	"database/sql"
	"errors"

	"github.com/doug-martin/goqu/v9"
)

type Location struct {
	ID      int    `db:"id" json:"id"`
	Street  string `db:"street" json:"street"`
	Number  string `db:"number" json:"number"`
	City    string `db:"city" json:"city"`
	Country string `db:"country" json:"country"`
}

const locationTable = "location"

func (s *Storage) LocationGetAll() ([]Location, error) {
	var locations = []Location{}
	sqlStr, _, err := goqu.Dialect("mysql").
		Select("*").
		From(locationTable).ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(sqlStr)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var location Location
		columns := getColumnsForStruct(&location)
		if err := rows.Scan(columns...); err != nil {
			return nil, err
		}
		locations = append(locations, location)
	}

	return locations, nil
}

func (s *Storage) Location(locationID int) (*Location, error) {
	var location = &Location{}
	sqlStr, _, err := goqu.Dialect(s.dialect).
		From(locationTable).
		Select("*").
		Where(goqu.C("id").Eq(locationID)).ToSQL()
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRow(sqlStr)
	if err != nil {
		return nil, err
	}

	columns := getColumnsForStruct(location)
	err = row.Scan(columns...)
	if err != nil {
		return nil, err
	}

	return location, nil
}

func (s *Storage) InsertLocation(location *Location) (int64, error) {
	sqlStr, _, err := goqu.Dialect(s.dialect).
		From(locationTable).
		Insert().
		Rows(location).ToSQL()
	if err != nil {
		return 0, err
	}

	result, err := s.db.Exec(sqlStr)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (s *Storage) UpdateLocation(location *Location) (*Location, error) {
	updatelocation := &Location{}
	sqlStr, _, err := goqu.Dialect(s.dialect).
		From(locationTable).
		Update().
		Set(location).
		Where(goqu.C("id").Eq(location.ID)).ToSQL()
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRow(sqlStr)
	if err != nil {
		return nil, err
	}

	if row.Err() != nil {
		return nil, row.Err()
	}

	columns := getColumnsForStruct(updatelocation)
	err = row.Scan(columns...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return location, nil
		}
		return nil, err
	}

	return location, nil
}

func (s *Storage) DeleteLocation(locationID int) (*Location, error) {
	sqlStr, _, err := goqu.Dialect(s.dialect).
		From(locationTable).
		Delete().
		Where(goqu.C("id").Eq(locationID)).ToSQL()
	if err != nil {
		return nil, err
	}

	location, err := s.Location(locationID)
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(sqlStr)
	if err != nil {
		return nil, err
	}

	return location, nil
}
