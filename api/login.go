package api

import (
	"errors"
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

//func Login(secret string) func(w http.ResponseWriter, r *http.Request) {
//	return func(w http.ResponseWriter, r *http.Request) {
//		ctx := r.Context()
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
//		lr := &loginRequest{}
//		if err := json.NewDecoder(r.Body).Decode(lr); err != nil {
//			log.Error(ctx, herrors.New("Error decoding incoming request", "req", string(body)))
//			ResponseJSONError(ctx, w, "Invalid payload", http.StatusInternalServerError)
//			return
//		}
//
//		// Validate
//		if err := lr.validate(); err != nil {
//			log.Info(ctx, herrors.New("Invalid login data"))
//			ResponseJSONError(ctx, w, err.Error(), http.StatusBadRequest)
//			return
//		}
//
//		// We pull the user and then filter if we really need to do the hashing, since that might be an expensive
//		// operation!?
//		var userid int64
//		var hashedPassword []byte
//		q := `SELECT id, password FROM users WHERE email = $1;`
//		err = db.Instance.QueryRowContext(ctx, q, lr.Email).Scan(&userid, &hashedPassword)
//		if err != nil && err != sql.ErrNoRows {
//			log.Critical(ctx, herrors.Wrap(err))
//			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
//			return
//		}
//
//		if err == sql.ErrNoRows || userid == 0 {
//			log.Info(ctx, herrors.New("User not found", "email", lr.Email))
//			ResponseJSONError(ctx, w, "Username/password incorrect", http.StatusBadRequest)
//			return
//		}
//
//		if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(lr.Password)); err != nil {
//			log.Error(ctx, herrors.Wrap(err))
//			ResponseJSONError(ctx, w, "Username/password incorrect", http.StatusBadRequest)
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
