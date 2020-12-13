package api

import (
	"encoding/json"
	"net/http"

	"github.com/satori/go.uuid"

	"github.com/nazarnovak/hobee-be/pkg/db"
	"github.com/nazarnovak/hobee-be/pkg/herrors2"
	"github.com/nazarnovak/hobee-be/pkg/log"
)

type ContactRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func (cr *ContactRequest) Validate() error {
	if cr.Name == "" {
		return herrors.New("Please provide your name")
	}

	if cr.Email == "" {
		return herrors.New("Please provide your email")
	}

	if cr.Message == "" {
		return herrors.New("Please provide your message")
	}

	if !emailValidationRegEx.MatchString(cr.Email) {
		return herrors.New("Please provide a valid email")
	}

	return nil
}

func Contact(secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		//if err := checkOrigin(r); err != nil {
		//	log.Critical(ctx, err)
		//	ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
		//	return
		//}

		// TODO: If user already logged in, save who was it?
		uuidStr, err := getCookieUUID(r, secret)
		if err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}
		//
		//if uuidStr == "" {
		//	log.Critical(ctx, herrors.New("Attempting to access messages without being logged in"))
		//	ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
		//	return
		//}

		cr := ContactRequest{}

		if err := json.NewDecoder(r.Body).Decode(&cr); err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		if err := cr.Validate(); err != nil {
			ResponseJSONError(ctx, w, err.Error(), http.StatusBadRequest)
			return
		}

		// Save feedback entry into DB
		if err := SaveContact(cr, uuidStr); err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		responseJSONSuccess(ctx, w)
	}
}

func SaveContact(cr ContactRequest, userUUID string) error {
	q := `INSERT INTO contacts(id, name, email, message, useruuid)
		VALUES(DEFAULT, $1, $2, $3, $4);`

	var err error
	u := uuid.Nil
	if userUUID != "" {
		u, err = uuid.FromString(userUUID)
		if err != nil {
			return herrors.Wrap(err)
		}
	}

	if _, err := db.Instance.Exec(q, cr.Name, cr.Email, cr.Message, u); err != nil {
		return herrors.Wrap(err)
	}

	return nil
}
