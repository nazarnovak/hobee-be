package api

import (
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

func loggedInUserId(r *http.Request, secret string) string {
	keyFunc := func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	}

	cookie, _ := r.Cookie(sessionCookieName)
	if cookie == nil {
		// Log.Error
		println("loggedInUserId 1")
		return ""
	}

	tkn, err := jwt.Parse(cookie.Value, keyFunc)
	if err != nil {
		// Log.Error
		println("loggedInUserId 2:", err.Error())
		return ""
	}

	claims, ok := tkn.Claims.(jwt.MapClaims)
	if !ok {
		// Log.Error
		println("loggedInUserId 3")
		return ""
	}

	userid, ok := claims["userid"]
	if !ok {
		// Log.Error
		println("loggedInUserId 4")
		return ""
	}

	useridStr, ok := userid.(string)
	if !ok {
		// Log.Error
		println("loggedInUserId 5")
		return ""
	}

	return useridStr
}

func User(secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userIdStr := loggedInUserId(r, secret)

		if userIdStr == "" {
			if err := responseJSONError(w, "Not authorized", http.StatusBadRequest); err != nil {
				// Log.Error
				println("user 1:", err.Error())
			}
		}

		s := struct {
			UserID string `json:"userid"`
		}{
			UserID: userIdStr,
		}

		if err := responseJSONObject(w, s); err != nil {
			// Log.Error
			println("user 2:", err.Error())
		}
	}
}
