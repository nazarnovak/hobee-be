package controllers

import (
	"net/http"
)

type inputEmailLogin struct {
	Email    string
	Password string
	Test     bool
}

func Login(w http.ResponseWriter, r *http.Request) {
	//ctx := r.Context()


}
