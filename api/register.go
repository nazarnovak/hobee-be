package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"hobee-be/pkg/log"
	"net/http"
	"regexp"
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

func ResponseJSONError(ctx context.Context, w http.ResponseWriter, msg string, status int) {
	jr := jsonResponse{Error: true, Msg: msg}

	b, err := json.Marshal(jr)
	if err != nil {
		log.Critical(ctx, err)
		return
	}

	http.Error(w, string(b), status)
}

func responseJSONSuccess(ctx context.Context, w http.ResponseWriter) {
	jr := jsonResponse{Error: false, Msg: "Success"}

	b, err := json.Marshal(jr)
	if err != nil {
		log.Critical(ctx, err)
		return
	}

	http.Error(w, string(b), http.StatusOK)
}

func responseJSONObject(ctx context.Context, w http.ResponseWriter, obj interface{}) {
	b, err := json.Marshal(obj)
	if err != nil {
		log.Critical(ctx, err)
		return
	}

	http.Error(w, string(b), http.StatusOK)
}

//func Register(secret string) func(w http.ResponseWriter, r *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		ctx := r.Context()
//		defer r.Body.Close()
//
//		if r == nil {
//			log.Critical(ctx, herrors.New("request is nil"))
//			return
//		}
//
//		body, err := ioutil.ReadAll(r.Body)
//		if err != nil {
//			log.Critical(ctx, herrors.New("Could not ReadAll from Body"))
//			return
//		}
//
//		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
//
//		if userIdStr := LoggedInUserId(r, secret); userIdStr != "" {
//			log.Info(ctx, herrors.New("Already logged in"))
//			ResponseJSONError(ctx, w, "Already logged in", http.StatusInternalServerError)
//			return
//		}
//
//		rr := &registerRequest{}
//		if err := json.NewDecoder(r.Body).Decode(rr); err != nil {
//			log.Error(ctx, herrors.New("Error decoding incoming request", "req", string(body)))
//			ResponseJSONError(ctx, w, "Invalid payload", http.StatusInternalServerError)
//			return
//		}
//
//		// Nothing to sanitize?
//		// Validate
//		if err := rr.validate(); err != nil {
//			log.Info(ctx, herrors.New("Invalid register data"))
//			ResponseJSONError(ctx, w, err.Error(), http.StatusBadRequest)
//			return
//		}
//
//		// Check if email already taken
//		var exists bool
//		q := `SELECT 1 FROM users WHERE email = $1;`
//		err = db.Instance.QueryRowContext(ctx, q, rr.Email).Scan(&exists)
//		if err != nil && err != sql.ErrNoRows {
//			log.Critical(ctx, herrors.Wrap(err))
//			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
//			return
//		}
//		if exists {
//			log.Info(ctx, herrors.New("Email already taken", "email", rr.Email))
//			ResponseJSONError(ctx, w, "Email already taken", http.StatusBadRequest)
//			return
//		}
//
//		// Check if invitationcode exists
//		var invitationCodeId, max int
//		q = `SELECT id, max FROM invitationcodes WHERE code = $1;`
//		err = db.Instance.QueryRowContext(ctx, q, rr.InvitationCode).Scan(&invitationCodeId, &max)
//		if err != nil && err != sql.ErrNoRows {
//			log.Critical(ctx, herrors.Wrap(err))
//			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
//			return
//		}
//		if invitationCodeId == 0 {
//			log.Warning(ctx, herrors.New("Invitation code not found", "code", rr.InvitationCode))
//			ResponseJSONError(ctx, w, "Invitation code not found", http.StatusBadRequest)
//			return
//		}
//
//		// Check if invitationcode limit reached
//		var usersWithInvitationCount int
//		q = `SELECT COUNT(*) FROM users WHERE invitationcodeid = $1;`
//		err = db.Instance.QueryRowContext(ctx, q, invitationCodeId).Scan(&usersWithInvitationCount)
//		if err != nil && err != sql.ErrNoRows {
//			log.Critical(ctx, herrors.Wrap(err))
//			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
//			return
//		}
//		if usersWithInvitationCount >= max {
//			log.Warning(ctx, herrors.New("Invitation code limit reached"))
//			ResponseJSONError(ctx, w, "Invitation code limit reached", http.StatusForbidden)
//			return
//		}
//
//		byteHashedPassword, err := bcrypt.GenerateFromPassword([]byte(rr.Password), bcrypt.DefaultCost)
//		if err != nil {
//			log.Error(ctx, herrors.Wrap(err))
//			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
//			return
//		}
//
//		var userid int64
//		q = `INSERT INTO users(id, email, password, invitationcodeid, created) VALUES(DEFAULT, $1, $2, $3, DEFAULT) returning id;`
//		if err := db.Instance.QueryRowContext(ctx, q, rr.Email, byteHashedPassword, invitationCodeId).Scan(&userid); err != nil {
//			log.Critical(ctx, herrors.Wrap(err))
//			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
//			return
//		}
//
//		if userid == 0 {
//			log.Critical(ctx, herrors.New("Could not insert a new user", "email", rr.Email))
//			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
//			return
//		}
//
//		// JWT + cookie
//		claims := jwt.MapClaims{
//			"userid": strconv.FormatInt(userid, 10),
//		}
//
//		tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
//		signed, err := tkn.SignedString([]byte(secret))
//		if err != nil {
//			log.Critical(ctx, herrors.Wrap(err))
//			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
//			return
//		}
//
//		c := &http.Cookie{
//			Path:   "/",
//			Name:   sessionCookieName,
//			Value:  signed,
//			MaxAge: sessionTimeInSeconds,
//		}
//
//		http.SetCookie(w, c)
//
//		responseJSONSuccess(ctx, w)
//	}
//}
