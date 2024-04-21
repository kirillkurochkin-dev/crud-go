package rest

import (
	"context"
	"crud-go/internal/entity"
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strconv"
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

	return r
}

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
