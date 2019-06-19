package api

import (
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"

	"github.com/nazarnovak/hobee-be/pkg/log"
	"github.com/nazarnovak/hobee-be/pkg/herrors"
)

func TestLogin(secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		keys, ok := r.URL.Query()["id"]

		if !ok || len(keys[0]) < 1 {
			log.Info(ctx, herrors.New("Url Param 'id' is missing"))
			return
		}

		id, err := strconv.ParseInt(keys[0], 10, 64)
		if err != nil {
			log.Error(ctx, herrors.Wrap(err))
			return
		}

		// JWT + cookie
		claims := jwt.MapClaims{
			"userid": strconv.FormatInt(id, 10),
		}

		tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := tkn.SignedString([]byte(secret))
		if err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		c := &http.Cookie{
			Path:   "/",
			Name:   sessionCookieName,
			Value:  signed,
			MaxAge: sessionTimeInSeconds,
		}

		http.SetCookie(w, c)

		responseJSONSuccess(ctx, w)
	}
}

func TestLogout(secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// JWT + cookie
		claims := jwt.MapClaims{
			"userid": strconv.FormatInt(0, 10),
		}

		tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := tkn.SignedString([]byte(secret))
		if err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		c := &http.Cookie{
			Path:   "/",
			Name:   sessionCookieName,
			Value:  signed,
			MaxAge: 0,
		}

		http.SetCookie(w, c)

		responseJSONSuccess(ctx, w)
	}
}
