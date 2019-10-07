package actions

import (
	"encoding/json"
	"net/http"

	"github.com/nlopes/slack"
	"github.com/omaressameldin/lunch-roulette/pkg/commands"
)

type Reply struct {
	Attachments []slack.Attachment `json:"attachments"`
	Text        string             `json:"text"`
}

func unmarshalPayload(r *http.Request) (*slack.InteractionCallback, error) {
	var payload slack.InteractionCallback
	err := json.Unmarshal([]byte(r.FormValue("payload")), &payload)
	if err != nil {
		return nil, err
	}

	return &payload, nil
}

func mightCancel(payload *slack.InteractionCallback, w http.ResponseWriter) bool {
	isCancel := payload.ActionCallback.AttachmentActions[0].Name == commands.CancelValue
	if isCancel {
		sendCancelResponse(w, canceledRequest)
	}

	return isCancel
}

func sendCancelResponse(w http.ResponseWriter, text string) {
	sendReply(w, Reply{
		Attachments: []slack.Attachment{commands.DangerMessage(text)},
	})
}

func sendReply(w http.ResponseWriter, r Reply) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r)
}
