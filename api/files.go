package api

import (
	"fmt"
	"net/http"

	"hobee-be/pkg/hconst"
)

func Files(w http.ResponseWriter, r *http.Request) error {
	file := fmt.Sprintf("%s/%s", hconst.PublicFolder, r.URL.Path)

	http.ServeFile(w, r, file)

	return nil
}
