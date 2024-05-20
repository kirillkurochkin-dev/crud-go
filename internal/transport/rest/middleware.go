package rest

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

type CtxValue int

const (
	ctxUserID CtxValue = iota
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"URI":    r.RequestURI,
		}).Info()
		next.ServeHTTP(w, r)
	})
}

func (c *Controller) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := getTokenFromRequest(r)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"handler": "authMiddleware",
				"problem": "getTokenFromRequest error",
			}).Error(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		userId, err := c.usersService.ParseToken(r.Context(), token)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"handler": "authMiddleware",
				"problem": "service error",
			}).Error(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), ctxUserID, userId)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func getTokenFromRequest(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")
	if header == "" {
		return "", errors.New("empty auth header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid auth header")
	}

	if len(headerParts[1]) == 0 {
		return "", errors.New("token is empty")
	}

	return headerParts[1], nil
}
