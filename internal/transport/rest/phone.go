package rest

import (
	"context"
	"crud-go/internal/entity"
	"encoding/json"
	"io"
	"net/http"

	_ "crud-go/docs"

	"github.com/sirupsen/logrus"
)

// @Summary Get a phone by ID
// @Description Retrieve a phone record by its ID
// @Tags Phones
// @Accept json
// @Produce json
// @Param id path int true "Phone ID"
// @Success 200 {object} entity.Phone "OK"
// @Failure 400 {string} string "Bad Request"
// @Router /api/phones/{id} [get]
func (c *Controller) getPhoneById(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromReq(r)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "getPhoneById",
			"problem": "getting id from request",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	book, err := c.phonesService.GetPhoneById(context.TODO(), id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "getPhoneById",
			"problem": "service error",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	response, err := json.Marshal(book)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "getPhoneById",
			"problem": "marshal error",
		}).Error(err)
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
func (c *Controller) getAllPhones(w http.ResponseWriter, r *http.Request) {
	phones, err := c.phonesService.GetAllPhones(context.TODO())
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "getAllPhones",
			"problem": "service error",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(phones)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "getAllPhones",
			"problem": "marshal error",
		}).Error(err)
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
func (c *Controller) createPhone(w http.ResponseWriter, r *http.Request) {
	var phone entity.PhoneInputDto

	reqBytes, err := io.ReadAll(r.Body)

	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "createPhone",
			"problem": "reading body",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBytes, &phone)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "createPhone",
			"problem": "unmarshal error",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.phonesService.CreatePhone(context.TODO(), phone)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "createPhone",
			"problem": "service error",
		}).Error(err)
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
func (c *Controller) updatePhoneById(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromReq(r)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "updatePhoneById",
			"problem": "getting id from request",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var phone entity.PhoneInputDto

	reqBytes, err := io.ReadAll(r.Body)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "updatePhoneById",
			"problem": "reading body",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = json.Unmarshal(reqBytes, &phone)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "updatePhoneById",
			"problem": "unmarshal error",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c.phonesService.UpdatePhoneById(context.TODO(), id, phone)
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
func (c *Controller) deletePhoneById(w http.ResponseWriter, r *http.Request) {
	id, err := getIdFromReq(r)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "deletePhoneById",
			"problem": "getting id from request",
		}).Error(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = c.phonesService.DeletePhoneById(context.TODO(), id)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "deletePhoneById",
			"problem": "service error",
		}).Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
