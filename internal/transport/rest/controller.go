package rest

import (
	"context"
	"crud-go/internal/entity"
	"errors"
	"net/http"
	"strconv"

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

type UsersService interface {
	SignUp(ctx context.Context, input entity.SignUpInput) error
	SignIn(ctx context.Context, input entity.SignInInput) (string, error)
	ParseToken(ctx context.Context, token string) (int64, error)
}

type Controller struct {
	phonesService PhonesService
	usersService  UsersService
}

func NewController(phonesService PhonesService, usersService UsersService) *Controller {
	return &Controller{
		phonesService: phonesService,
		usersService:  usersService,
	}
}

func (c *Controller) InitRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(loggingMiddleware)

	auth := r.PathPrefix("/api/users").Subrouter()
	{
		auth.HandleFunc("/sign-up", c.signUp).Methods(http.MethodPost)
		auth.HandleFunc("/sign-in", c.signIn).Methods(http.MethodPost)
	}

	phones := r.PathPrefix("/api/phones").Subrouter()
	{
		phones.Use(c.authMiddleware)
		phones.HandleFunc("", c.createPhone).Methods(http.MethodPost)
		phones.HandleFunc("", c.getAllPhones).Methods(http.MethodGet)
		phones.HandleFunc("/{id:[0-9]+", c.getPhoneById).Methods(http.MethodGet)
		phones.HandleFunc("/{id:[0-9]+", c.deletePhoneById).Methods(http.MethodDelete)
		phones.HandleFunc("/{id:[0-9]+", c.updatePhoneById).Methods(http.MethodPut)
	}

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"), //The url pointing to API definition

	))

	return r
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
