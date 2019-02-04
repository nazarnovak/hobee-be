/*
	Request to "/"
*/

package controllers

import (
	"fmt"
	"html/template"
	"net/http"

	"hobee-be/pkg/hconst"
)

func Index(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles(fmt.Sprintf("%s/index.gohtml", hconst.ViewFolder)))
	tpl.Execute(w, nil)
}
