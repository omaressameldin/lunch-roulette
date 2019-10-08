package actions

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/nlopes/slack"
	"github.com/omaressameldin/lunch-roulette/pkg/commands"
)

type Reply struct {
	Attachments []slack.Attachment `json:"attachments"`
	Blocks      []slack.Block      `json:"blocks"`
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
	isCancel := payload.ActionCallback.BlockActions[0].Value == commands.CancelValue
	if isCancel {
		sendCancelResponse(payload.ResponseURL, w, canceledRequest)
	}

	return isCancel
}

func sendCancelResponse(url string, w http.ResponseWriter, text string) {
	sendReply(url, w, Reply{
		Attachments: []slack.Attachment{commands.DangerMessage(text)},
	})
}

func sendPendingResponse(url string, w http.ResponseWriter) {
	sendReply(url, w, Reply{
		Attachments: []slack.Attachment{commands.PendingResponse()},
	})
}

func sendReply(url string, w http.ResponseWriter, r Reply) {
	jsonValue, _ := json.Marshal(r)
	log.Println(string(jsonValue))
	http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
}
