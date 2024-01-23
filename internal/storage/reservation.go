package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/doug-martin/goqu/v9"
)

var reservationTable = "reservation"

type Reservation struct {
	ID          int    `db:"id"`
	ContactName string `db:"contactName"`
	PhoneNumber string `db:"phoneNumber"`
	HolidayID   int    `db:"holidayID"`
}

type ReservationResult struct {
	ID          int                 `db:"id" json:"id"`
	ContactName string              `db:"contactName" json:"contactName"`
	PhoneNumber string              `db:"phoneNumber" json:"phoneNumber"`
	Holiday     HolidayWithLocation `db:"holiday" json:"holiday"`
}

func (s *Storage) ReservationGetAll() (interface{}, error) {
	sqlStr, _, err := goqu.Dialect("mysql").
		Select(goqu.T(reservationTable).All(), goqu.T(holidaysTable).All(), goqu.T(locationTable).All()).
		From(reservationTable).InnerJoin(
		goqu.T(holidaysTable),
		goqu.On(goqu.Ex{reservationTable + ".holidayID": goqu.I(holidaysTable + ".id")}),
	).InnerJoin(
		goqu.T(locationTable),
		goqu.On(goqu.Ex{holidaysTable + ".locationID": goqu.I(locationTable + ".id")}),
	).ToSQL()
	if err != nil {
		return nil, err
	}

	fmt.Println(sqlStr)

	rows, err := s.db.Query(sqlStr)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	resultStruct := []ReservationResult{}

	for rows.Next() {
		var reservation Reservation
		var holiday Holiday
		var location Location
		columns := getColumnsForStruct(&reservation)
		columns = append(columns, getColumnsForStruct(&holiday)...)
		columns = append(columns, getColumnsForStruct(&location)...)
		if err := rows.Scan(columns...); err != nil {
			return nil, err
		}
		resultStruct = append(resultStruct, ReservationResult{
			ID:          reservation.ID,
			ContactName: reservation.ContactName,
			PhoneNumber: reservation.PhoneNumber,
			Holiday: HolidayWithLocation{
				ID:        holiday.ID,
				Title:     holiday.Title,
				Duration:  holiday.Duration,
				StartDate: holiday.StartDate,
				Price:     holiday.Price,
				FreeSlots: holiday.FreeSlots,
				Location:  location,
			},
		})
	}

	return resultStruct, nil
}

func (s *Storage) Reservation(reservationID int) (*Reservation, error) {
	var reservation = &Reservation{}
	sqlStr, _, err := goqu.Dialect(s.dialect).
		From(reservationTable).
		Select("*").
		Where(goqu.C("id").Eq(reservationID)).ToSQL()
	if err != nil {
		return nil, err
	}

	row := s.db.QueryRow(sqlStr)
	if err != nil {
		return nil, err
	}

	columns := getColumnsForStruct(reservation)
	err = row.Scan(columns...)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}

func (s *Storage) InsertReservation(reservation *Reservation) (int64, error) {
	sqlStr, _, err := goqu.Dialect(s.dialect).
		From(reservationTable).
		Insert().
		Rows(reservation).ToSQL()
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

func (s *Storage) UpdateReservation(reservation *Reservation) (*Reservation, error) {
	updateReservation := &Reservation{}
	sqlStr, _, err := goqu.Dialect(s.dialect).
		From(reservationTable).
		Update().
		Set(reservation).
		Where(goqu.C("id").Eq(reservation.ID)).ToSQL()
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

	columns := getColumnsForStruct(updateReservation)
	err = row.Scan(columns...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return reservation, nil
		}
		return nil, err
	}

	return reservation, nil
}

func (s *Storage) DeleteReservation(reservationID int) (*Reservation, error) {
	sqlStr, _, err := goqu.Dialect(s.dialect).
		From(reservationTable).
		Delete().
		Where(goqu.C("id").Eq(reservationID)).ToSQL()
	if err != nil {
		return nil, err
	}

	reservation, err := s.Reservation(reservationID)
	if err != nil {
		return nil, err
	}

	_, err = s.db.Exec(sqlStr)
	if err != nil {
		return nil, err
	}

	return reservation, nil
}
