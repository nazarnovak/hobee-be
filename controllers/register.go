package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"hobee-be/pkg/db"
)

const (
	emailMaxLenght    = 255
	passwordMinLength = 8
	passwordMaxLenght = 64
)

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

	// Check if email already taken
	var exists bool
	q := `SELECT 1 FROM users WHERE email = $1;`
	err := db.Instance.QueryRowContext(ctx, q, rr.Email).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		// Log.Error
		fmt.Println("Email collision error:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if exists {
		err := errors.New("Email already taken")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if invitationcode exists
	var id, max int
	q = `SELECT id, max FROM invitationcodes WHERE code = $1;`
	err = db.Instance.QueryRowContext(ctx, q, rr.InvitationCode).Scan(&id, &max)
	if err != nil && err != sql.ErrNoRows {
		// Log.Error
		fmt.Println("Invitation collision error:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if max == 0 {
		err := errors.New("Invitation code not found")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if invitationcode limit reached
	var usersWithInvitationCount int
	q = `SELECT COUNT(*) FROM users WHERE invitationcodeid = $1;`
	err = db.Instance.QueryRowContext(ctx, q, id).Scan(&usersWithInvitationCount)
	if err != nil && err != sql.ErrNoRows {
		// Log.Error
		fmt.Println("Invitation collision error:", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if usersWithInvitationCount >= max {
		err := errors.New("Invitation code limit reached")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Set password
	// add a cookie in the response!
}
