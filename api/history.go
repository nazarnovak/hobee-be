package api

import (
	"encoding/csv"
	"fmt"
	"github.com/nazarnovak/hobee-be/pkg/socket"
	"net/http"
	"os"
	"time"

	"github.com/nazarnovak/hobee-be/pkg/herrors2"
	"github.com/nazarnovak/hobee-be/pkg/log"
)

type HistoryResponse struct {
	Rooms map[string][]socket.Message `json:"rooms"`
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

		roomHistory, err := socket.UserRoomHistory(uuidStr)
		if err != nil {
			log.Critical(ctx, err)
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		o := HistoryResponse{
			Rooms: map[string][]socket.Message{},
		}

		errorCount := 0
		for _, roomUUID := range roomHistory {
			// Open the csv file and get the messages
			file, err := os.OpenFile(fmt.Sprintf("%s/%s.csv", "chats", roomUUID), os.O_RDONLY, 0777)
			if err != nil {
				errorCount++
				continue
			}
			defer file.Close()

			csvReader := csv.NewReader(file)
			csvReader.Comma = ';'
			csvReader.LazyQuotes = true

			records, err := csvReader.ReadAll()
			if err != nil {
				log.Critical(ctx, herrors.Wrap(err))
				continue
			}

			// Skip header line
			records = append(records[:0], records[1:]...)

			messages, err := recordsToMessages(records)
			if err != nil {
				log.Critical(ctx, herrors.Wrap(err, "roomuuid", roomUUID))
				continue
			}

			//o.Rooms[roomUUID] = make([]socket.Message, 0, len(messages))
			o.Rooms[roomUUID] = messages
		}

		if errorCount > 1 {
			log.Critical(ctx, herrors.New("Couldn't open some files to view chat history", "roomuuids",
				roomHistory))
		}

		// All 3 files couldn't be opened, we just return internal error
		if errorCount == 3 {
			ResponseJSONError(ctx, w, internalServerError, http.StatusInternalServerError)
			return
		}

		responseJSONObject(ctx, w, o)
	}
}

func recordsToMessages(records [][]string) ([]socket.Message, error) {
	msgs := []socket.Message{}

	for _, record := range records {
		ts, err := time.Parse(time.RFC3339, record[0])
		if err != nil {
			return nil, herrors.Wrap(err)
		}

		msg := socket.Message{
			AuthorUUID: record[1],
			Type: socket.MessageType(record[2]),
			Text: record[3],
			Timestamp: ts,
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}
