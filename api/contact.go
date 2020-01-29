package api

import (
	"encoding/json"
	"github.com/nazarnovak/hobee-be/pkg/herrors2"
	"net/http"
	//"github.com/nazarnovak/hobee-be/pkg/email"
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
		//uuidStr, err := getCookieUUID(r, secret)
		//if err != nil {
		//	log.Critical(ctx, herrors.Wrap(err))
		//	ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
		//	return
		//}
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

		//subject := "New feedback"
		//text := fmt.Sprintf("Name: %s\nEmail: %s\nMessage: %s\n", cr.Name, cr.Email, cr.Message)
		//if err := email.Send(subject, text); err != nil {
		//	log.Critical(ctx, herrors.Wrap(err))
		//	ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
		//	return
		//}

		responseJSONSuccess(ctx, w)
	}
}
