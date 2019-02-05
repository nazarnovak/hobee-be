package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

const (
	emailMaxLenght = 255
	passwordMinLength = 8
	passwordMaxLenght = 64
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	InvitationCode string `json:"invitationCode"`
}

func (rr *registerRequest) validate() error {
	if len(rr.Email) == 0 {
		return errors.New("Email cannot be empty")
	}

	if len(rr.Email) > emailMaxLenght {
		return errors.New(fmt.Sprintf("Email cannot be longer than %d characters", emailMaxLenght))
	}

	if len(rr.Password) == 0 {
		return errors.New("Email cannot be empty")
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

func Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	defer r.Body.Close()

	rr := &registerRequest{}
	if err := json.NewDecoder(r.Body).Decode(rr); err != nil {
		// Log.Error
		http.Error(w, responseSomethingWentWrong, http.StatusInternalServerError)
		return
	}

	// Nothing to sanitize?
	// Validate
	if err := rr.validate(); err != nil {
		// Log.Error
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check email collision
	// check if invitation code exists and if we didn't reach a limit on invitation codes
	// add a cookie in the response!
}
