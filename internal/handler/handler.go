package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
	"travel/internal/service"

	"github.com/gorilla/mux"
)

type Service interface {
	ReservationGetAll() (interface{}, error)
	Reservation(reservationID int) (*service.ReservationDTO, error)
	InsertReservation(reservation service.ReservationDTO) (int64, error)
	UpdateReservation(reservation service.ReservationDTO) (*service.ReservationDTO, error)
	DeleteReservation(reservationID int) (*service.ReservationDTO, error)

	LocationGetAll() ([]service.LocationDTO, error)
	Location(locationID int) (*service.LocationDTO, error)
	InsertLocation(Location service.LocationDTO) (int64, error)
	UpdateLocation(Location service.LocationDTO) (*service.LocationDTO, error)
	DeleteLocation(locationID int) (*service.LocationDTO, error)

	HolidayGetAll(filterDTO service.FilterHolidays) (interface{}, error)
	Holiday(holidayID int) (*service.HolidayDTO, error)
	InsertHoliday(Holiday service.HolidayDTO) (int64, error)
	UpdateHoliday(Holiday service.HolidayDTO) (*service.HolidayDTO, error)
	DeleteHoliday(holidayID int) (*service.HolidayDTO, error)
}

type apiHandler struct {
	service Service
}

func New(service Service) http.Handler {
	handler := &apiHandler{service: service}

	//create route
	route := mux.NewRouter()

	//holidays
	route.Methods(http.MethodGet).Path("/holidays").HandlerFunc(handler.GetHolidays)
	route.Methods(http.MethodGet).Path("/holidays/{id}").HandlerFunc(handler.GetHoliday)
	route.Methods(http.MethodPost).Path("/holidays").HandlerFunc(handler.CreateHoliday)
	route.Methods(http.MethodPut).Path("/holidays").HandlerFunc(handler.UpdateHoliday)
	route.Methods(http.MethodDelete).Path("/holidays/{id}").HandlerFunc(handler.DeleteHoliday)

	//locations
	route.Methods(http.MethodGet).Path("/locations").HandlerFunc(handler.GetLocations)
	route.Methods(http.MethodGet).Path("/locations/{id}").HandlerFunc(handler.GetLocation)
	route.Methods(http.MethodPost).Path("/locations").HandlerFunc(handler.CreateLocation)
	route.Methods(http.MethodPut).Path("/locations").HandlerFunc(handler.UpdateLocation)
	route.Methods(http.MethodDelete).Path("/locations/{id}").HandlerFunc(handler.DeleteLocation)

	//reservations
	route.Methods(http.MethodGet).Path("/reservations").HandlerFunc(handler.GetReservations)
	route.Methods(http.MethodGet).Path("/reservations/{id}").HandlerFunc(handler.GetReservation)
	route.Methods(http.MethodPost).Path("/reservations").HandlerFunc(handler.CreateReservation)
	route.Methods(http.MethodPut).Path("/reservations").HandlerFunc(handler.UpdateReservation)
	route.Methods(http.MethodDelete).Path("/reservations/{id}").HandlerFunc(handler.DeleteReservation)

	return route
}

func (h *apiHandler) GetHolidays(w http.ResponseWriter, r *http.Request) {
	location := r.FormValue("location")

	var duration int
	var startDate time.Time
	var err error

	if strings.TrimSpace(r.FormValue("duration")) != "" {
		duration, err = strconv.Atoi(r.FormValue("duration"))
		if err != nil {
			jsonResponseWrite(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	if strings.TrimSpace(r.FormValue("startDate")) != "" {
		startDate, err = time.Parse(time.DateOnly, r.FormValue("startDate"))
		if err != nil {
			jsonResponseWrite(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	holidays, err := h.service.HolidayGetAll(service.FilterHolidays{
		StartDate: startDate,
		Duration:  duration,
		Location:  location,
	})

	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponseWrite(w, holidays, http.StatusOK)
}

func (h *apiHandler) GetHoliday(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusBadRequest)
		return
	}

	holiday, err := h.service.Holiday(id)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponseWrite(w, holiday, http.StatusOK)
}

func (h *apiHandler) CreateHoliday(w http.ResponseWriter, r *http.Request) {
	data := RequestHoliday{}

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		fmt.Println("error: ", err)
		jsonResponseWrite(w, err.Error(), http.StatusInternalServerError)
		return
	}

	startDate, err := time.Parse(time.DateOnly, data.StartDate)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusBadRequest)
		return
	}

	price, err := strconv.ParseFloat(data.Price, 64)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusBadRequest)
		return
	}

	idResult, err := h.service.InsertHoliday(service.HolidayDTO{
		Title:      data.Title,
		StartDate:  startDate,
		Duration:   data.Duration,
		Price:      float64(price),
		FreeSlots:  data.FreeSlots,
		LocationID: data.LocationID,
	})

	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponseWrite(w, idResult, http.StatusOK)
}

func (h *apiHandler) UpdateHoliday(w http.ResponseWriter, r *http.Request) {
	holiday := service.HolidayDTO{}

	err := json.NewDecoder(r.Body).Decode(&holiday)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusInternalServerError)
		return
	}

	idResult, err := h.service.UpdateHoliday(holiday)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponseWrite(w, idResult, http.StatusOK)
}

// recive the id only
func (h *apiHandler) DeleteHoliday(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusBadRequest)
		return
	}

	holiday, err := h.service.DeleteHoliday(id)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponseWrite(w, holiday, http.StatusOK)
}

func (h *apiHandler) GetLocations(w http.ResponseWriter, r *http.Request) {
	locations, err := h.service.LocationGetAll()
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponseWrite(w, locations, http.StatusOK)
}

// recive the id only
func (h *apiHandler) GetLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusBadRequest)
		return
	}

	location, err := h.service.Location(id)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponseWrite(w, location, http.StatusOK)
}

func (h *apiHandler) CreateLocation(w http.ResponseWriter, r *http.Request) {
	location := service.LocationDTO{}

	err := json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusInternalServerError)
		return
	}

	idResult, err := h.service.InsertLocation(location)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponseWrite(w, idResult, http.StatusOK)
}

func (h *apiHandler) UpdateLocation(w http.ResponseWriter, r *http.Request) {
	location := service.LocationDTO{}

	err := json.NewDecoder(r.Body).Decode(&location)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusInternalServerError)
		return
	}

	updatedLocation, err := h.service.UpdateLocation(location)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponseWrite(w, updatedLocation, http.StatusOK)
}

// recive the id only
func (h *apiHandler) DeleteLocation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusBadRequest)
		return
	}

	location, err := h.service.DeleteLocation(id)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponseWrite(w, location, http.StatusOK)
}

func (h *apiHandler) GetReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := h.service.ReservationGetAll()
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponseWrite(w, reservations, http.StatusOK)
}

// recive the id only
func (h *apiHandler) GetReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusBadRequest)
		return
	}

	reservation, err := h.service.Reservation(id)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponseWrite(w, reservation, http.StatusOK)
}

func (h *apiHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	reservation := service.ReservationDTO{}

	err := json.NewDecoder(r.Body).Decode(&reservation)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusInternalServerError)
		return
	}

	idResult, err := h.service.InsertReservation(reservation)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponseWrite(w, idResult, http.StatusOK)
}

func (h *apiHandler) UpdateReservation(w http.ResponseWriter, r *http.Request) {
	reservation := service.ReservationDTO{}

	err := json.NewDecoder(r.Body).Decode(&reservation)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusInternalServerError)
		return
	}

	result, err := h.service.UpdateReservation(reservation)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsonResponseWrite(w, result, http.StatusOK)
}

// recive the id only
func (h *apiHandler) DeleteReservation(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusBadRequest)
		return
	}

	reservations, err := h.service.DeleteReservation(id)
	if err != nil {
		jsonResponseWrite(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsonResponseWrite(w, reservations, http.StatusOK)
}

func jsonResponseWrite(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
