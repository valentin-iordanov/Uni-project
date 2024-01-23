package service

import (
	"fmt"
	"time"
	"travel/internal/storage"
)

type Storage interface {
	//reservation
	ReservationGetAll() (interface{}, error)
	Reservation(reservationID int) (*storage.Reservation, error)
	InsertReservation(reservation *storage.Reservation) (int64, error)
	UpdateReservation(reservation *storage.Reservation) (*storage.Reservation, error)
	DeleteReservation(reservationID int) (*storage.Reservation, error)

	//location
	LocationGetAll() ([]storage.Location, error)
	Location(locationID int) (*storage.Location, error)
	InsertLocation(location *storage.Location) (int64, error)
	UpdateLocation(location *storage.Location) (*storage.Location, error)
	DeleteLocation(locationID int) (*storage.Location, error)

	//holiday
	HolidaysGetAll(location string, duration int, startDate time.Time) ([]storage.HolidayWithLocation, error)
	Holiday(holidaysID int) (*storage.Holiday, error)
	InsertHolidays(holidays *storage.Holiday) (int64, error)
	UpdateHolidays(holidays *storage.Holiday) (*storage.Holiday, error)
	DeleteHolidays(holidaysID int) (*storage.Holiday, error)
}

type Service struct {
	storage Storage
}

func New(storage Storage) *Service {
	return &Service{storage: storage}
}

func (s *Service) ReservationGetAll() (interface{}, error) {
	reservations, err := s.storage.ReservationGetAll()
	if err != nil {
		return nil, err
	}

	return reservations, nil
}

func (s *Service) Reservation(reservationID int) (*ReservationDTO, error) {
	reservation, err := s.storage.Reservation(reservationID)
	if err != nil {
		return nil, err
	}

	result := &ReservationDTO{
		ID:          reservation.ID,
		ContactName: reservation.ContactName,
		PhoneNumber: reservation.PhoneNumber,
		HolidayID:   reservation.HolidayID,
	}

	return result, nil

}

func (s *Service) InsertReservation(reservation ReservationDTO) (int64, error) {

	reservationData := &storage.Reservation{
		ID:          reservation.ID,
		ContactName: reservation.ContactName,
		PhoneNumber: reservation.PhoneNumber,
		HolidayID:   reservation.HolidayID,
	}

	return s.storage.InsertReservation(reservationData)
}

func (s *Service) UpdateReservation(reservation ReservationDTO) (*ReservationDTO, error) {
	reservationData := &storage.Reservation{
		ID:          reservation.ID,
		ContactName: reservation.ContactName,
		PhoneNumber: reservation.PhoneNumber,
		HolidayID:   reservation.HolidayID,
	}

	updatedReservation, err := s.storage.UpdateReservation(reservationData)
	if err != nil {
		return nil, err
	}

	reservation = ReservationDTO{
		ID:          updatedReservation.ID,
		ContactName: updatedReservation.ContactName,
		PhoneNumber: updatedReservation.PhoneNumber,
		HolidayID:   updatedReservation.HolidayID,
	}

	return &reservation, nil
}

func (s *Service) DeleteReservation(reservationID int) (*ReservationDTO, error) {
	reservation, err := s.storage.DeleteReservation(reservationID)
	if err != nil {
		return nil, err
	}

	result := &ReservationDTO{
		ID:          reservation.ID,
		ContactName: reservation.ContactName,
		PhoneNumber: reservation.PhoneNumber,
		HolidayID:   reservation.HolidayID,
	}

	return result, nil
}

func (s *Service) LocationGetAll() ([]LocationDTO, error) {
	locations, err := s.storage.LocationGetAll()
	if err != nil {
		return nil, err
	}

	result := []LocationDTO{}
	for _, value := range locations {
		result = append(result, LocationDTO{
			ID:      value.ID,
			Street:  value.Street,
			Number:  value.Number,
			City:    value.City,
			Country: value.Country,
		})
	}

	return result, nil
}

func (s *Service) Location(locationID int) (*LocationDTO, error) {
	location, err := s.storage.Location(locationID)
	if err != nil {
		return nil, err
	}

	result := &LocationDTO{
		ID:      location.ID,
		Street:  location.Street,
		Number:  location.Number,
		City:    location.City,
		Country: location.Country,
	}

	return result, nil

}

func (s *Service) InsertLocation(location LocationDTO) (int64, error) {

	locationData := &storage.Location{
		ID:      location.ID,
		Street:  location.Street,
		Number:  location.Number,
		City:    location.City,
		Country: location.Country,
	}

	return s.storage.InsertLocation(locationData)
}

func (s *Service) UpdateLocation(location LocationDTO) (*LocationDTO, error) {
	reservationData := &storage.Location{
		ID:      location.ID,
		Street:  location.Street,
		Number:  location.Number,
		City:    location.City,
		Country: location.Country,
	}

	updatedReservation, err := s.storage.UpdateLocation(reservationData)
	if err != nil {
		return nil, err
	}

	location = LocationDTO{
		ID:      updatedReservation.ID,
		Street:  updatedReservation.Street,
		Number:  updatedReservation.Number,
		City:    updatedReservation.City,
		Country: updatedReservation.Country,
	}

	return &location, nil
}

func (s *Service) DeleteLocation(locationID int) (*LocationDTO, error) {
	location, err := s.storage.DeleteLocation(locationID)
	if err != nil {
		return nil, err
	}

	result := &LocationDTO{
		ID:      location.ID,
		Street:  location.Street,
		Number:  location.Number,
		City:    location.City,
		Country: location.Country,
	}

	return result, nil
}

func (s *Service) HolidayGetAll(filterHolidays FilterHolidays) (interface{}, error) {
	holidays, err := s.storage.HolidaysGetAll(filterHolidays.Location, filterHolidays.Duration, filterHolidays.StartDate)
	if err != nil {
		return nil, err
	}

	return holidays, nil
}

func (s *Service) Holiday(holidayID int) (*HolidayDTO, error) {
	holiday, err := s.storage.Holiday(holidayID)
	if err != nil {
		return nil, err
	}

	result := &HolidayDTO{
		ID:         holiday.ID,
		Title:      holiday.Title,
		StartDate:  holiday.StartDate,
		Duration:   holiday.Duration,
		Price:      holiday.Price,
		FreeSlots:  holiday.FreeSlots,
		LocationID: holiday.LocationID,
	}

	return result, nil

}

func (s *Service) InsertHoliday(holiday HolidayDTO) (int64, error) {

	holidayData := &storage.Holiday{
		Title:      holiday.Title,
		StartDate:  holiday.StartDate,
		Duration:   holiday.Duration,
		Price:      holiday.Price,
		FreeSlots:  holiday.FreeSlots,
		LocationID: holiday.LocationID,
	}

	fmt.Printf("holidayData: %v\n", holidayData)

	return s.storage.InsertHolidays(holidayData)
}

func (s *Service) UpdateHoliday(holiday HolidayDTO) (*HolidayDTO, error) {
	reservationData := &storage.Holiday{
		ID:         holiday.ID,
		Title:      holiday.Title,
		StartDate:  holiday.StartDate,
		Duration:   holiday.Duration,
		Price:      holiday.Price,
		FreeSlots:  holiday.FreeSlots,
		LocationID: holiday.LocationID,
	}

	updatedReservation, err := s.storage.UpdateHolidays(reservationData)
	if err != nil {
		return nil, err
	}

	holiday = HolidayDTO{
		ID:         updatedReservation.ID,
		Title:      holiday.Title,
		StartDate:  holiday.StartDate,
		Duration:   holiday.Duration,
		Price:      holiday.Price,
		FreeSlots:  holiday.FreeSlots,
		LocationID: holiday.LocationID,
	}

	return &holiday, nil
}

func (s *Service) DeleteHoliday(holidayID int) (*HolidayDTO, error) {
	holiday, err := s.storage.DeleteHolidays(holidayID)
	if err != nil {
		return nil, err
	}

	result := &HolidayDTO{
		ID:         holiday.ID,
		Title:      holiday.Title,
		StartDate:  holiday.StartDate,
		Duration:   holiday.Duration,
		Price:      holiday.Price,
		FreeSlots:  holiday.FreeSlots,
		LocationID: holiday.LocationID,
	}

	return result, nil
}
