package api

import (
"fmt"
	"encoding/json"
	"net/http"
	"time"

	"github.com/satori/go.uuid"

	"github.com/nazarnovak/hobee-be/pkg/db"
	"github.com/nazarnovak/hobee-be/pkg/herrors2"
	"github.com/nazarnovak/hobee-be/pkg/log"
)

type FeedbackRequest struct {
	Message string `json:"message"`
}

func (fr *FeedbackRequest) Validate() error {
	if fr.Message == "" {
		return herrors.New("Please provide your message")
	}

	return nil
}

func Feedback(secret string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		//if err := checkOrigin(r); err != nil {
		//	log.Critical(ctx, err)
		//	ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
		//	return
		//}
fmt.Println("here")
		// TODO: If user already logged in, save who was it?
		uuidStr, err := getCookieUUID(r, secret)
		if err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}
		
		if uuidStr == "" {
			log.Critical(ctx, herrors.New("Attempting to access feedback without being logged in"))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		fr := FeedbackRequest{}

		if err := json.NewDecoder(r.Body).Decode(&fr); err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		if err := fr.Validate(); err != nil {
			ResponseJSONError(ctx, w, err.Error(), http.StatusBadRequest)
			return
		}

		// Save feedback entry into DB
		if err := SaveFeedback(fr.Message, uuidStr); err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		responseJSONSuccess(ctx, w)
	}
}

func SaveFeedback(message, userUUID string) error {
	q := `INSERT INTO feedbacks(id, message, useruuid, created)
		VALUES(DEFAULT, $1, $2, $3);`

	var err error

	u := uuid.Nil
	if userUUID != "" {
		u, err = uuid.FromString(userUUID)
		if err != nil {
			return herrors.Wrap(err)
		}
	}

	if _, err := db.Instance.Exec(q, message, u, time.Now().UTC()); err != nil {
		return herrors.Wrap(err)
	}

	return nil
}
