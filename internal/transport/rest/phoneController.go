package rest

import (
	"context"
	"crud-go/internal/entity"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	_ "crud-go/docs"

	"github.com/gorilla/mux"

	httpSwagger "github.com/swaggo/http-swagger"
)

type PhonesService interface {
	GetPhoneById(ctx context.Context, id int64) (entity.Phone, error)
	GetAllPhones(ctx context.Context) ([]entity.Phone, error)
	CreatePhone(ctx context.Context, ph entity.PhoneInputDto) error
	UpdatePhoneById(ctx context.Context, id int64, ph entity.PhoneInputDto) error
	DeletePhoneById(ctx context.Context, id int64) error
}

type Phones struct {
	phonesService PhonesService
}

func NewPhonesHandler(phonesService PhonesService) *Phones {
	return &Phones{
		phonesService: phonesService,
	}
}

func (p *Phones) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)
	phones := r.PathPrefix("/api/phones").Subrouter()
	{
		phones.HandleFunc("", p.createPhone).Methods(http.MethodPost)
		phones.HandleFunc("", p.getAllPhones).Methods(http.MethodGet)
		phones.HandleFunc("/{id:[0-9]+", p.getPhoneById).Methods(http.MethodGet)
		phones.HandleFunc("/{id:[0-9]+", p.deletePhoneById).Methods(http.MethodDelete)
		phones.HandleFunc("/{id:[0-9]+", p.updatePhoneById).Methods(http.MethodPut)
	}

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition

	))

	return r
}

// @Summary Get a phone by ID
// @Description Retrieve a phone record by its ID
// @Tags Phones
// @Accept json
// @Produce json
// @Param id path int true "Phone ID"
// @Success 200 {object} entity.Phone "OK"
// @Failure 400 {string} string "Bad Request"
// @Router /api/phones/{id} [get]
func (p *Phones) getPhoneById(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromReq(r)
	if err != nil {
		log.Println("GetPhoneById() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := p.phonesService.GetPhoneById(context.TODO(), id)
	if err != nil {
		log.Println("GetPhoneById() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(book)
	if err != nil {
		log.Println("GetPhoneById() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Get all phones
// @Description Retrieve all phone records
// @Tags Phones
// @Accept json
// @Produce json
// @Success 200 {array} entity.Phone "OK"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/phones [get]
func (p *Phones) getAllPhones(w http.ResponseWriter, r *http.Request) {
	phones, err := p.phonesService.GetAllPhones(context.TODO())
	if err != nil {
		log.Println("GetAllPhones() error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(phones)
	if err != nil {
		log.Println("GetAllPhones() error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response)
}

// @Summary Create a new phone
// @Description Create a new phone record
// @Tags Phones
// @Accept json
// @Produce json
// @Param phone body entity.PhoneInputDto true "Phone Data"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/phones [post]
func (p *Phones) createPhone(w http.ResponseWriter, r *http.Request) {
	var phone entity.PhoneInputDto

	reqBytes, err := io.ReadAll(r.Body)

	if err != nil {
		log.Println("CreatePhone() error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBytes, &phone)
	if err != nil {
		log.Println("CreatePhone() error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = p.phonesService.CreatePhone(context.TODO(), phone)
	if err != nil {
		log.Println("CreatePhone() error: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// @Summary Update a phone by ID
// @Description Update an existing phone record
// @Tags Phones
// @Accept json
// @Produce json
// @Param id path int true "Phone ID"
// @Param phone body entity.PhoneInputDto true "Phone Data"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Router /api/phones/{id} [put]neInputDto true "Phone Data"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Router /api/phones/{id} [put]
func (p *Phones) updatePhoneById(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromReq(r)
	if err != nil {
		log.Println("UpdatePhoneById() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var phone entity.PhoneInputDto

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("UpdatePhoneById() error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBytes, &phone)
	if err != nil {
		log.Println("CreatePhone() error: ", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	p.phonesService.UpdatePhoneById(context.TODO(), id, phone)
	w.WriteHeader(http.StatusOK)
}

// @Summary Delete a phone by ID
// @Description Delete a phone record by its ID
// @Tags Phones
// @Accept json
// @Produce json
// @Param id path int true "Phone ID"
// @Success 200 {string} string "OK"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/phones/{id} [delete]
func (p *Phones) deletePhoneById(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromReq(r)
	if err != nil {
		log.Println("DeletePhoneById() error:", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = p.phonesService.DeletePhoneById(context.TODO(), id)
	if err != nil {
		log.Println("DeletePhoneById() error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getIdFromReq(r *http.Request) (int64, error) {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return 0, err
	}

	if id == 0 {
		return 0, errors.New("id can't be 0")
	}

	return id, nil
}
