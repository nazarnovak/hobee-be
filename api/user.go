package api

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"

	"github.com/nazarnovak/hobee-be/pkg/herrors2"
	"github.com/nazarnovak/hobee-be/pkg/log"
)

type Response struct {
	Issued bool `json:"issued"`
}

func isLoggedIn(r *http.Request, secret string) (bool, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}

	cookie, _ := r.Cookie(sessionCookieName)
	if cookie == nil {
		return false, nil
	}

	tkn, err := jwt.Parse(cookie.Value, keyFunc)
	if err != nil {
		return false, herrors.Wrap(err)
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return false, herrors.New("Could not assert token claims as jwt.MapClaims")
	}

	userAgent, ok := claims["user-agent"]
	if !ok {
		return false, herrors.New("Could not find user-agent in claims")
	}

	userAgentStr, ok := userAgent.(string)
	if !ok {
		return false, herrors.New("Could not assert user-agent as a string")
	}

	ip, ok := claims["ip"]
	if !ok {
		return false, herrors.New("Could not find ip in claims")
	}

	ipStr, ok := ip.(string)
	if !ok {
		return false, herrors.New("Could not assert ip as a string")
	}

	id, ok := claims["uuid"]
	if !ok {
		return false, herrors.New("Could not find uuid in claims")
	}

	idStr, ok := id.(string)
	if !ok {
		return false, herrors.New("Could not assert id as a string")
	}

	if userAgentStr == "" {
		return false, herrors.New("User-agent claim can not be empty")
	}

	if ipStr == "" {
		return false, herrors.New("IP claim can not be empty")
	}

	if idStr == "" {
		return false, herrors.New("id claim can not be empty")
	}

	return true, nil
}

func getCookieUUID(r *http.Request, secret string) (string, error) {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}

	cookie, _ := r.Cookie(sessionCookieName)
	if cookie == nil {
		return "", nil
	}

	tkn, err := jwt.Parse(cookie.Value, keyFunc)
	if err != nil {
		return "", herrors.Wrap(err)
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		return "", herrors.New("Could not assert token claims as jwt.MapClaims")
	}

	id, ok := claims["uuid"]
	if !ok {
		return "", herrors.New("Could not find uuid in claims")
	}

	idStr, ok := id.(string)
	if !ok {
		return "", herrors.New("Could not assert id as a string")
	}

	return idStr, nil
}

// Identify checks if the request contains a logged in cookie or we need to set one
func Identify(secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		loggedIn, err := isLoggedIn(r, secret)
		if err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		if loggedIn {
			responseJSONSuccess(ctx, w)
			return
		}

		id := uuid.NewV4()

		// JWT + cookie
		claims := jwt.MapClaims{
			"ip":         r.RemoteAddr,
			"user-agent": r.UserAgent(),
			"uuid":       id.String(),
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
