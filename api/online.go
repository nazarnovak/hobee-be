package api

import (
	"net/http"

	"github.com/nazarnovak/hobee-be/pkg/herrors2"
	"github.com/nazarnovak/hobee-be/pkg/log"
	"github.com/nazarnovak/hobee-be/pkg/socket"
)

type onlineResponse struct {
	Total        int `json:"total"`
	Talking      int `json:"talking"`
	Searching    int `json:"searching"`
	Disconnected int `json:"disconnected"`
}

func Online(secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		//if err := checkOrigin(r); err != nil {
		//	log.Critical(ctx, err)
		//	ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
		//	return
		//}
		//
		//uuidStr, err := getCookieUUID(r, secret)
		//if err != nil {
		//	log.Critical(ctx, herrors.Wrap(err))
		//	ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
		//	return
		//}
		//
		//if uuidStr == "" {
		//	log.Critical(ctx, herrors.New("Attempting to access websockets without being logged in"))
		//	ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
		//	return
		//}

		user, pass, ok := r.BasicAuth()
		if !ok {
			log.Critical(ctx, herrors.New("No basic auth to  access to online"))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		if user != "n" || pass != "n" {
			log.Critical(ctx, herrors.New("Basic auth wrong credentials", "user", user, "pass", pass))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		talking, searching, disconnected := socket.GetTalkingSearchingDisconnected()

		o := onlineResponse{
			Total:        socket.GetTotalOnline(),
			Talking:      talking,
			Searching:    searching,
			Disconnected: disconnected,
		}

		responseJSONObject(ctx, w, o)
	}
}
