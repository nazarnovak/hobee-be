package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"regexp"
	"strconv"

	"golang.org/x/crypto/bcrypt"
	"hobee-be/pkg/db"
)

const (
	internalServerError = "Internal server error"

	emailMaxLenght    = 255
	passwordMinLength = 8
	passwordMaxLenght = 64
)

var emailValidationRegEx = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@([a-zA-Z0-9-]+\\.)+[a-zA-Z0-9-]{2,}$")

type registerRequest struct {
	Email          string `json:"email"`
	Password       string `json:"password"`
	InvitationCode string `json:"invitationCode"`
}

type invitationCode struct {
	Id   int
	Code string
	Max  int
}

func (rr *registerRequest) validate() error {
	if len(rr.Email) == 0 {
		return errors.New("Email cannot be empty")
	}

	if len(rr.Email) > emailMaxLenght {
		return errors.New(fmt.Sprintf("Email cannot be longer than %d characters", emailMaxLenght))
	}

	if !emailValidationRegEx.MatchString(rr.Email) {
		return errors.New("Email format incorrect")
	}

	if len(rr.Password) == 0 {
		return errors.New("Passowrd cannot be empty")
	}

	if len(rr.Password) < passwordMinLength {
		return errors.New(fmt.Sprintf("Password cannot be shorter than %d characters", passwordMinLength))
	}

	if len(rr.Password) > passwordMaxLenght {
		return errors.New(fmt.Sprintf("Password cannot be longer than %d characters", passwordMaxLenght))
	}

	if len(rr.InvitationCode) == 0 {
		return errors.New("Invitation code cannot be empty")
	}

	return nil
}

type jsonResponse struct {
	Error bool   `json:"error"`
	Msg   string `json:"msg"`
}

func responseJSONError(w http.ResponseWriter, msg string, status int) error {
	jr := jsonResponse{Error: true, Msg: msg}

	b, err := json.Marshal(jr)
	if err != nil {
		return err
	}

	http.Error(w, string(b), status)

	return nil
}

func responseJSONSuccess(w http.ResponseWriter) error {
	jr := jsonResponse{Error: false, Msg: "Success"}

	b, err := json.Marshal(jr)
	if err != nil {
		return err
	}

	http.Error(w, string(b), http.StatusOK)

	return nil
}

func responseJSONObject(w http.ResponseWriter, obj interface{}) error {
	b, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	http.Error(w, string(b), http.StatusOK)

	return nil
}

func Register(secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		defer r.Body.Close()

		// TODO: Respond with codes to FE? So if you're already logged in I can return code 1 and then map that on
		// the FE? And if it's 1 I can redirect somewhere else
		if userIdStr := loggedInUserId(r, secret); userIdStr != "" {
			println("register -1")
			if err := responseJSONError(w, "Already logged in", http.StatusInternalServerError); err != nil {
				// Log.Error
				println("register 0")
			}
			return
		}

		rr := &registerRequest{}
		if err := json.NewDecoder(r.Body).Decode(rr); err != nil {
			// Log.Error
			println("register 1:", err.Error())
			if err = responseJSONError(w, "Invalid payload", http.StatusInternalServerError); err != nil {
				// Log.Error
				println("register 2:", err.Error())
			}
			return
		}

		// Nothing to sanitize?
		// Validate
		if err := rr.validate(); err != nil {
			// Log.Error
			println("register 3:", err.Error())
			if err = responseJSONError(w, err.Error(), http.StatusBadRequest); err != nil {
				// Log.Error
				println("register 4:", err.Error())
			}
			return
		}

		// Check if email already taken
		var exists bool
		q := `SELECT 1 FROM users WHERE email = $1;`
		err := db.Instance.QueryRowContext(ctx, q, rr.Email).Scan(&exists)
		if err != nil && err != sql.ErrNoRows {
			// Log.Error
			println("register 5:", err.Error())
			if err = responseJSONError(w, internalServerError, http.StatusInternalServerError); err != nil {
				// Log.Error
				println("register 6:", err.Error())
			}
			return
		}
		if exists {
			if err = responseJSONError(w, "Email already taken", http.StatusBadRequest); err != nil {
				// Log.Error
				println("register 7:", err.Error())
			}
			return
		}

		// Check if invitationcode exists
		var invitationCodeId, max int
		q = `SELECT id, max FROM invitationcodes WHERE code = $1;`
		err = db.Instance.QueryRowContext(ctx, q, rr.InvitationCode).Scan(&invitationCodeId, &max)
		if err != nil && err != sql.ErrNoRows {
			// Log.Error
			println("register 8:", err.Error())
			if err = responseJSONError(w, internalServerError, http.StatusInternalServerError); err != nil {
				// Log.Error
				println("register 9:", err.Error())
			}
			return
		}
		if max == 0 {
			if err = responseJSONError(w, "Invitation code not found", http.StatusBadRequest); err != nil {
				// Log.Error
				println("register 10:", err.Error())
			}
			return
		}

		// Check if invitationcode limit reached
		var usersWithInvitationCount int
		q = `SELECT COUNT(*) FROM users WHERE invitationcodeid = $1;`
		err = db.Instance.QueryRowContext(ctx, q, invitationCodeId).Scan(&usersWithInvitationCount)
		if err != nil && err != sql.ErrNoRows {
			// Log.Error
			println("register 11:", err.Error())
			if err = responseJSONError(w, internalServerError, http.StatusInternalServerError); err != nil {
				// Log.Error
				println("register 12:", err.Error())
			}
			return
		}
		if usersWithInvitationCount >= max {
			if err = responseJSONError(w, "Invitation code limit reached", http.StatusForbidden); err != nil {
				// Log.Error
				println("register 13:", err.Error())
			}
			return
		}

		byteHashedPassword, err := bcrypt.GenerateFromPassword([]byte(rr.Password), bcrypt.DefaultCost)
		if err != nil {
			// Log.Error
			println("register 13.1:", err.Error())
			if err = responseJSONError(w, internalServerError, http.StatusInternalServerError); err != nil {
				// Log.Error
				println("register 13.2:", err.Error())
			}
			return
		}

		var userid int64
		q = `INSERT INTO users(id, email, password, invitationcodeid, created) VALUES(DEFAULT, $1, $2, $3, DEFAULT) returning id;`
		if err := db.Instance.QueryRowContext(ctx, q, rr.Email, byteHashedPassword, invitationCodeId).Scan(&userid);
		err != nil {
			// Log.Error
			println("register 14:", err.Error())
			if err = responseJSONError(w, internalServerError, http.StatusInternalServerError); err != nil {
				// Log.Error
				println("register 15:", err.Error())
			}
			return
		}

		if userid == 0 {
			// Log.Error
			println("register 16")
			if err = responseJSONError(w, internalServerError, http.StatusInternalServerError); err != nil {
				// Log.Error
				println("register 17:", err.Error())
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
			println("register 18:", err.Error())
			if err = responseJSONError(w, internalServerError, http.StatusInternalServerError); err != nil {
				// Log.Error
				println("register 19:", err.Error())
			}
			return
		}

		c := &http.Cookie{
			Name: sessionCookieName,
			Value: signed,
			MaxAge: 600,
		}

		http.SetCookie(w, c)

		if err := responseJSONSuccess(w); err != nil {
			// Log.Error
			println("register 20:", err.Error())
			return
		}
	}
}
