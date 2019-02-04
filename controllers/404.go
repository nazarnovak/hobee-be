/*
	404 page
*/

package controllers

import (
	"fmt"
	"html/template"
	"net/http"

	"hobee-be/pkg/hconst"
)

func GetNotFound(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		println("Not get")
		return
	}

	tpl := template.Must(template.ParseFiles(fmt.Sprintf("%s/404.gohtml", hconst.ViewFolder)))
	tpl.Execute(w, nil)
}
