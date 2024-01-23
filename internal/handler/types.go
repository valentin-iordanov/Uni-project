package handler

type RequestHoliday struct {
	ID         string `json:"id"`
	Title      string `json:"title"`
	Duration   int    `json:"duration"`
	StartDate  string `json:"startDate"`
	Price      string `json:"price"`
	FreeSlots  int    `json:"freeSlots"`
	LocationID int    `json:"location"`
}
