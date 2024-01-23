package storage

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/doug-martin/goqu/v9"
)

type Holiday struct {
	ID         int       `db:"id"`
	Title      string    `db:"title"`
	Duration   int       `db:"duration"`
	StartDate  time.Time `db:"startDate"`
	Price      float64   `db:"price"`
	FreeSlots  int       `db:"freeSlots"`
	LocationID int       `db:"locationID"`
}

type HolidayWithLocation struct {
	ID        int       `db:"id" json:"id"`
	Title     string    `db:"title" json:"title"`
	Duration  int       `db:"duration" json:"duration"`
	StartDate time.Time `db:"startDate" json:"startDate"`
	Price     float64   `db:"price" json:"price"`
	FreeSlots int       `db:"freeSlots" json:"freeSlots"`
	Location  Location  `json:"location"`
}

const holidaysTable = "holiday"

func (s *Storage) HolidaysGetAll(location string, duration int, startDate time.Time) ([]HolidayWithLocation, error) {
	var holidays = []HolidayWithLocation{}
	sql := goqu.Dialect("mysql").
		Select(goqu.T(holidaysTable).All(), goqu.T(locationTable).All()).
		From(holidaysTable)

	sql = sql.InnerJoin(
		goqu.T(locationTable),
		goqu.On(goqu.Ex{holidaysTable + ".locationID": goqu.I(locationTable + ".id")}),
	)

	if location != "" || duration > 0 || !startDate.IsZero() {
		if location != "" {
			sql = sql.Where(goqu.ExOr{
				locationTable + ".country": goqu.Op{"like": "%" + location + "%"},
				locationTable + ".city":    goqu.Op{"like": "%" + location + "%"},
			})
		}

		if duration > 0 {
			sql = sql.Where(goqu.ExOr{
				"duration": goqu.Op{"eq": duration},
			})
		}

		if !startDate.IsZero() {
			sql = sql.Where(goqu.ExOr{
				"startDate": goqu.Op{"eq": startDate},
			})
		}
	}
	sqlStr, _, err := sql.ToSQL()
	if err != nil {
		return nil, err
	}

	rows, err := s.db.Query(sqlStr)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var holiday Holiday
		columns := getColumnsForStruct(&holiday)
		var location Location
		columns = append(columns, getColumnsForStruct(&location)...)
		if err := rows.Scan(columns...); err != nil {
			fmt.Printf("holidays: %v\n", holidays)
			fmt.Printf("err: %v\n", err.Error())
			return nil, err
		}
		holidays = append(holidays, HolidayWithLocation{
			ID:        holiday.ID,
			Title:     holiday.Title,
			Duration:  holiday.Duration,
			StartDate: holiday.StartDate,
			Price:     holiday.Price,
			FreeSlots: holiday.FreeSlots,
			Location:  location,
		})
	}

	return holidays, nil
}

func (s *Storage) Holiday(holidaysID int) (*Holiday, error) {
	var holidays = &Holiday{}
	sqlStr, _, err := goqu.Dialect(s.dialect).
		From(holidaysTable).
		Select("*").
		Where(goqu.C("id").Eq(holidaysID)).ToSQL()
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRow(sqlStr)
	if err != nil {
		return nil, err
	}

	columns := getColumnsForStruct(holidays)
	err = row.Scan(columns...)
	if err != nil {
		return nil, err
	}

	return holidays, nil
}

func (s *Storage) InsertHolidays(holidays *Holiday) (int64, error) {
	sqlStr, _, err := goqu.Dialect(s.dialect).
		From(holidaysTable).
		Insert().
		Rows(holidays).ToSQL()
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

func (s *Storage) UpdateHolidays(holidays *Holiday) (*Holiday, error) {
	updateHolidays := &Holiday{}
	sqlStr, _, err := goqu.Dialect(s.dialect).
		From(holidaysTable).
		Update().
		Set(holidays).
		Where(goqu.C("id").Eq(holidays.ID)).ToSQL()
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

	columns := getColumnsForStruct(updateHolidays)
	err = row.Scan(columns...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return holidays, nil
		}
		return nil, err
	}

	return holidays, nil
}

func (s *Storage) DeleteHolidays(holidaysID int) (*Holiday, error) {
	sqlStr, _, err := goqu.Dialect(s.dialect).
		Delete(holidaysTable).
		Where(goqu.C("id").Eq(holidaysID)).ToSQL()
	if err != nil {
		return nil, err
	}

	holiday, err := s.Holiday(holidaysID)
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(sqlStr)
	if err != nil {
		return nil, err
	}

	return holiday, nil
}
