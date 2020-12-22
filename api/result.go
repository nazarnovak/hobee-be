package api

import (
	"net/http"

	"github.com/nazarnovak/hobee-be/pkg/herrors2"
	"github.com/nazarnovak/hobee-be/pkg/log"
	"github.com/nazarnovak/hobee-be/pkg/socket"
)

func Result(secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		//if err := checkOrigin(r); err != nil {
		//	log.Critical(ctx, err)
		//	ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
		//	return
		//}

		uuidStr, err := getCookieUUID(r, secret)
		if err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		if uuidStr == "" {
			log.Critical(ctx, herrors.New("Attempting to access result without being logged in"))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		roomUUID := socket.UserInARoomUUID(uuidStr)
		if roomUUID == "" {
			log.Critical(ctx, herrors.New("Attempting to pull result when user is not part of a room"))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		userResult, err := socket.GetRoomUserResults(roomUUID, uuidStr)
		if err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		// This is done because if we pull the data from the DB field, even if it's text[], it will still send as ´nil´over JSON
		// This makes sure JSON has an empty array instead, which simplifies things
		likes := []socket.LikeReason{}
		reports := []socket.ReportReason{}

		if userResult.Likes != nil {
			likes = userResult.Likes
		}

		if userResult.Reports != nil {
			reports = userResult.Reports
		}

		o := socket.Result {
			Likes: likes,
			Reports: reports,
		}

		responseJSONObject(ctx, w, o)
	}
}
