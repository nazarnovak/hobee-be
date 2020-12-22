package api

import (
	"encoding/csv"
	"fmt"
	"github.com/nazarnovak/hobee-be/pkg/socket"
	"net/http"
	"os"
	// "strconv"
	"time"

	"github.com/nazarnovak/hobee-be/pkg/herrors2"
	"github.com/nazarnovak/hobee-be/pkg/log"
)

type HistoryResponse struct {
	ChatsHistory []ChatHistory `json:"chats"`
}

type ChatHistory struct {
	Messages []socket.Message `json:"messages"`
	Result   socket.Result    `json:"result"`
	Duration string           `json:"duration"`
}

func History(secret string) func(w http.ResponseWriter, r *http.Request) {
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

		o := HistoryResponse{
			ChatsHistory: []ChatHistory{},
		}

		roomHistory, err := socket.UserRoomHistory(uuidStr)
		if err != nil {
			// User doesn't have chats yet, return empty object
			responseJSONObject(ctx, w, o)
			return
		}

		for _, roomUUID := range roomHistory {
			messages, err := getRoomMessages(roomUUID)
			if err != nil {
				log.Critical(ctx, err)
				continue
			}

			result, err := getRoomResult(roomUUID, uuidStr)
			if err != nil {
				log.Critical(ctx, err)
				continue
			}

			dur := getChatDuration(messages)

			c := ChatHistory{
				Messages: messages,
				Result: result,
				Duration: dur,
			}
			o.ChatsHistory = append(o.ChatsHistory, c)
		}

		responseJSONObject(ctx, w, o)
	}
}

func getRoomMessages(roomuuid string) ([]socket.Message, error) {
	filepath := fmt.Sprintf("%s/%s.csv", "chats", roomuuid)

	// Open the csv file and get the messages
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0777)
	if err != nil {
		return nil, herrors.Wrap(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'
	csvReader.LazyQuotes = true

	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, herrors.Wrap(err)
	}

	if len(rows) < 1 {
		return nil, herrors.New("Expecting at least 1 record in the csv", "roomuuid", roomuuid)
	}

	// Skip header line
	rows = append(rows[:0], rows[1:]...)

	messages, err := rowsToMessages(rows)
	if err != nil {
		return nil, herrors.Wrap(err)
	}

	return messages, nil
}

func rowsToMessages(rows [][]string) ([]socket.Message, error) {
	msgs := []socket.Message{}

	for _, row := range rows {
		ts, err := time.Parse(time.RFC3339, row[0])
		if err != nil {
			return nil, herrors.Wrap(err)
		}

		msg := socket.Message{
			AuthorUUID: row[1],
			Type:       socket.MessageType(row[2]),
			Text:       row[3],
			Timestamp:  ts,
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func getRoomResult(roomuuid, useruuid string) (socket.Result, error) {
	// Open the csv file and get the messages
	filepath := fmt.Sprintf("%s/%s:%s.csv", "chats", roomuuid, useruuid)

	// If results don't exist - it means the current user didn't like or report the conversation - return default values
	if exists := socket.FileExists(filepath); !exists {
		return socket.Result{}, nil
	}

	file, err := os.OpenFile(filepath, os.O_RDONLY, 0777)
	if err != nil {
		return socket.Result{}, herrors.Wrap(err)
	}
	defer file.Close()

	csvReader := csv.NewReader(file)
	csvReader.Comma = ';'
	csvReader.LazyQuotes = true

	rows, err := csvReader.ReadAll()
	if err != nil {
		return socket.Result{}, herrors.Wrap(err)
	}

	if len(rows) != 2 {
		return socket.Result{}, herrors.New("Expecting 2 rows in the csv", "roomuuid", roomuuid)
	}

	// liked, err := strconv.ParseBool(rows[1][0])
	// if err != nil {
	// 	return socket.Result{}, herrors.Wrap(err)
	// }

	r := socket.Result{
		// Liked:    liked,
		// Reported: socket.ReportReason(rows[1][1]),
	}

	return r, nil
}

func getChatDuration(msgs []socket.Message) string {
	if len(msgs) < 1 {
		return "0s"
	}

	l := len(msgs)

	// Subtract from the last message timestamp - first message timestamp
	d := msgs[l-1].Timestamp.Sub(msgs[0].Timestamp)

	if d.Hours() > 1 {
		return fmt.Sprintf("%dh %dm", int(d.Hours()), int(d.Minutes()))
	}

	return fmt.Sprintf("%dm %ds", int(d.Minutes()), int(d.Seconds()))
}
