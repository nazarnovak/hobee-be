package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"hobee-be/pkg/db"
)

type inputEmailLogin struct {
	Email    string
	Password string
	Test     bool
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (lr *loginRequest) validate() error {
	if len(lr.Email) == 0 {
		return errors.New("Email cannot be empty")
	}

	if !emailValidationRegEx.MatchString(lr.Email) {
		return errors.New("Email format incorrect")
	}

	if len(lr.Password) == 0 {
		return errors.New("Password cannot be empty")
	}

	return nil
}

func Login(secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		if userIdStr := loggedInUserId(r, secret); userIdStr != "" {
			println("login -1")
			if err := responseJSONError(w, "Already logged in", http.StatusInternalServerError); err != nil {
				// Log.Error
				println("login 0")
			}
			return
		}

		lr := &loginRequest{}
		if err := json.NewDecoder(r.Body).Decode(lr); err != nil {
			// Log.Error
			println("login 1:", err.Error())
			if err = responseJSONError(w, "Invalid payload", http.StatusInternalServerError); err != nil {
				// Log.Error
				println("login 2:", err.Error())
			}
			return
		}

		// Validate
		if err := lr.validate(); err != nil {
			// Log.Error
			println("login 3")
			if err = responseJSONError(w, err.Error(), http.StatusBadRequest); err != nil {
				// Log.Error
				println("login 4:", err.Error())
			}
			return
		}

		// We pull the user and then filter if we really need to do the hashing, since that might be an expensive
		// operation!?
		var userid int64
		var hashedPassword []byte
		q := `SELECT id, password FROM users WHERE email = $1;`
		err := db.Instance.QueryRowContext(ctx, q, lr.Email).Scan(&userid, &hashedPassword)
		if err != nil && err != sql.ErrNoRows {
			// Log.Error
			println("login 7:", err.Error())
			if err = responseJSONError(w, internalServerError, http.StatusInternalServerError); err != nil {
				// Log.Error
				println("login 8:", err.Error())
			}
			return
		}
println(userid)
		if err == sql.ErrNoRows || userid == 0 {
			// Log.Error
			println("login 8.1:", err.Error())
			if err = responseJSONError(w, "Username/password incorrect", http.StatusBadRequest); err != nil {
				// Log.Error
				println("login 9:", err.Error())
			}
			return
		}

		if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(lr.Password)); err != nil {
			// Log.Error
			println("login 9.1:", err.Error())
			if err = responseJSONError(w, "Username/password incorrect", http.StatusBadRequest); err != nil {
				// Log.Error
				println("login 9.2:", err.Error())
			}
			return
		}

		// JWT + cookie
		claims := jwt.MapClaims{
			"userid": strconv.FormatInt(userid, 10),
		}

		tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		signed, err := tkn.SignedString([]byte(secret))
		if err != nil {
			// Log.Error
			println("login 10:", err.Error())
			if err = responseJSONError(w, internalServerError, http.StatusInternalServerError); err != nil {
				// Log.Error
				println("login 11:", err.Error())
			}
			return
		}

		c := &http.Cookie{
			Name:   sessionCookieName,
			Value:  signed,
			MaxAge: 600,
		}

		http.SetCookie(w, c)

		if err := responseJSONSuccess(w); err != nil {
			// Log.Error
			println("login 12:", err.Error())
			return
		}
	}
}
