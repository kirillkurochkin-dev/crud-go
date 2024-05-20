package rest

import (
	"encoding/json"
	"io"
	"net/http"

	_ "crud-go/docs"
	"crud-go/internal/entity"

	"github.com/sirupsen/logrus"
)

// @Summary SignIn
// @Description SignIn
// @Tags Users
// @Accept json
// @Produce json
// @Param id path int true "Phone ID"
// @Success 200 {object} entity.Phone "OK"
// @Failure 400 {string} string "Bad Request"
// @Router /api/phones/{id} [get]
func (c *Controller) signIn(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "signIn",
			"problem": "reading body",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp entity.SignInInput
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "signIn",
			"problem": "unmarshal error",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := inp.Validate(); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "signUp",
			"problem": "validation error",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := c.usersService.SignIn(r.Context(), inp)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "signIn",
			"problem": "service error",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(map[string]string{
		"token": token,
	})
	if err != nil {
		logrus.WithFields(logrus.Fields{	
			"handler": "signIn",
			"problem": "marshal error",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// @Summary Sign up a new user]
// @Description Create a new user record
// @Tags Users
// @Accept json
// @Produce json
// @Param user body entity.SignUpInput true "User Data"
// @Success 201 {string} string "Created"
// @Failure 400 {string} string "Bad Request"
// @Failure 500 {string} string "Internal Server Error"
// @Router /api/users [post]
func (c *Controller) signUp(w http.ResponseWriter, r *http.Request) {
	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "signUp",
			"problem": "reading body",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var inp entity.SignUpInput
	if err = json.Unmarshal(reqBytes, &inp); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "signUp",
			"problem": "unmarshal error",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := inp.Validate(); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "signUp",
			"problem": "validation error",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.usersService.SignUp(r.Context(), inp)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "signUp",
			"problem": "service error",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
