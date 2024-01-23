package service

import "time"

type HolidayDTO struct {
	ID         int       `json:"id"`
	Title      string    `json:"title"`
	StartDate  time.Time `json:"startDate"`
	Duration   int       `json:"duration"`
	Price      float64   `json:"price"`
	FreeSlots  int       `json:"freeSlots"`
	LocationID int       `json:"location"`
}

type LocationDTO struct {
	ID      int    `json:"id"`
	Street  string `json:"street"`
	Number  string `json:"number"`
	City    string `json:"city"`
	Country string `json:"country"`
}

type FilterHolidays struct {
	Location  string
	StartDate time.Time
	Duration  int
}

type ReservationDTO struct {
	ID          int    `json:"id"`
	ContactName string `json:"contactName"`
	PhoneNumber string `json:"phoneNumber"`
	HolidayID   int    `json:"holiday"`
}
