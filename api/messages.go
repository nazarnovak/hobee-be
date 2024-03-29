package api

import (
	"net/http"

	"github.com/nazarnovak/hobee-be/pkg/herrors2"
	"github.com/nazarnovak/hobee-be/pkg/log"
	"github.com/nazarnovak/hobee-be/pkg/socket"
)

type MessagesResponse struct {
	Messages []socket.Message `json:"messages"`
}

func Messages(secret string) func(w http.ResponseWriter, r *http.Request) {
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
			log.Critical(ctx, herrors.New("Attempting to access messages without being logged in"))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		roomUUID := socket.UserInARoomUUID(uuidStr)
		if roomUUID == "" {
			log.Critical(ctx, herrors.New("Attempting to pull messages when user is not part of a room"))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		msgs, err := socket.RoomMessages(roomUUID)
		if err != nil {
			log.Critical(ctx, herrors.Wrap(err))
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		// We're marking users own messages so FE understands how to sort it, with removing uuids
		for k, msg := range msgs {
			if msg.Type == socket.MessageTypeSystem {
				msgs[k].AuthorUUID = string(socket.MessageTypeSystem)

				if msg.Text == socket.SystemDisconnected {
					msgs[k].AuthorUUID = string(socket.MessageTypeOwn)
					if msg.AuthorUUID != uuidStr {
						msgs[k].AuthorUUID = string(socket.MessageTypeBuddy)
					}
				}
			}

			if msg.Type == socket.MessageTypeActivity {
				msgs[k].AuthorUUID = string(socket.MessageTypeOwn)
				if msg.AuthorUUID != uuidStr {
					msgs[k].AuthorUUID = string(socket.MessageTypeBuddy)
				}
			}

			if msg.Type != socket.MessageTypeChatting {
				continue
			}

			msgs[k].AuthorUUID = string(socket.MessageTypeOwn)
			if msg.AuthorUUID != uuidStr {
				msgs[k].AuthorUUID = string(socket.MessageTypeBuddy)
				continue
			}
		}

		o := MessagesResponse{
			Messages: msgs,
		}

		responseJSONObject(ctx, w, o)
	}
}
