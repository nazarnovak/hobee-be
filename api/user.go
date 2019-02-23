package api

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"

	"hobee-be/pkg/herrors"
	"hobee-be/pkg/log"
)

func LoggedInUserId(r *http.Request, secret string) string {
	ctx := r.Context()

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}

	cookie, _ := r.Cookie(sessionCookieName)
	if cookie == nil {
		return ""
	}

	tkn, err := jwt.Parse(cookie.Value, keyFunc)
	if err != nil {
		log.Critical(ctx, herrors.Wrap(err))
		return ""
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		log.Critical(ctx, herrors.New("Could not assert token claims as jwt.MapClaims"))
		return ""
	}

	userid, ok := claims["userid"]
	if !ok {
		log.Critical(ctx, herrors.New("Could not find userid key in claims"))
		return ""
	}

	useridStr, ok := userid.(string)
	if !ok {
		log.Critical(ctx, herrors.New("Could not assert userid as a string"))
		return ""
	}

	// 0 is written in the claims when you log out
	if useridStr == "0" {
		return ""
	}

	return useridStr
}

func User(secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		userIdStr := LoggedInUserId(r, secret)
		if userIdStr == "" {
			ResponseJSONError(ctx, w, "Not authorized", http.StatusBadRequest)
		}

		s := struct {
			UserID string `json:"userid"`
		}{
			UserID: userIdStr,
		}

		responseJSONObject(ctx, w, s)
	}
}
