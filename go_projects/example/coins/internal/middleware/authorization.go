package middleware

import (
	"errors"
	"net/http"

	"github.com/HanmaDevin/example_coins/api"
	"github.com/HanmaDevin/example_coins/internal/tools"
	log "github.com/sirupsen/logrus"
)

var UnAutheriziedError = errors.New("Invalid username or token.")

func Authorization(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var username string = r.URL.Query().Get("username")
		var token = r.Header.Get("Authorization")
		var err error

		if username == "" || token == "" {
			log.Error(UnAutheriziedError)
			api.RequestErrorHandler(w, UnAutheriziedError)
			return
		}

		var database *tools.DatabaseInterface
		database, err = tools.NewDataBase()
		if err != nil {
			api.InternalErrorHandler(w)
			return
		}

		var loginDetails *tools.LoginDetails
		loginDetails = (*database).GetUserLoginDetails(username)

		if loginDetails == nil || (token != (*loginDetails).AuthToken) {
			log.Error(UnAutheriziedError)
			api.RequestErrorHandler(w, UnAutheriziedError)
			return
		}

		next.ServeHTTP(w, r)
	})
}
