package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/astapi/orchestra/aws/orchestra"
	"github.com/astapi/orchestra/aws/orchestra/ecs"
	"github.com/nlopes/slack"
)

// interactionHandler handles interactive message response.
type interactionHandler struct {
	slackClient       *slack.Client
	verificationToken string
}

func (h interactionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	if err != nil {
		log.Printf("[ERROR] Failed to unespace request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var message slack.AttachmentActionCallback
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Only accept message from slack with valid token
	if message.Token != h.verificationToken {
		log.Printf("[ERROR] Invalid token: %s", message.Token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	action := message.Actions[0]
	switch action.Name {
	case actionSelect:
		fmt.Println(action)
		value := action.SelectedOptions[0].Value

		text := ""
		if value == "deploy_all" {
			text = "デプロイ（全部入り）"
		}

		// Overwrite original drop down message.
		originalMessage := message.OriginalMessage
		originalMessage.Attachments[0].Text = fmt.Sprintf("%s を実行してよいですか？", text)
		originalMessage.Attachments[0].Actions = []slack.AttachmentAction{
			{
				Name:  actionStart,
				Text:  "Yes",
				Type:  "button",
				Value: "deploy_all",
				Style: "primary",
			},
			{
				Name:  actionCancel,
				Text:  "No",
				Type:  "button",
				Style: "danger",
			},
		}

		w.Header().Add("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(&originalMessage)
		return
	case actionStart:
		value := action.Value
		title := ":tada: deploy request done :+1:"

		sess, err := orchestra.CreateSession()
		if err != nil {
			title = err.Error()
			responseMessage(w, message.OriginalMessage, title, "")
			return
		}

		if value == "deploy_all" {
			// デプロイ実行（全部入り）
			ecs := ecs.New(sess)

			for task, cluster := range env.RunTasks {
				err := ecs.RunTask(cluster, task)
				if err != nil {
					title = err.Error()
					responseMessage(w, message.OriginalMessage, title, "")
					return
				}
			}

			for service, cluster := range env.DeployAllServices {
				err := ecs.UpdateService(cluster, service)
				if err != nil {
					title = err.Error()
					responseMessage(w, message.OriginalMessage, title, "")
					return
				}
			}
		}

		responseMessage(w, message.OriginalMessage, title, "")
		return
	case actionCancel:
		title := fmt.Sprintf(":x: @%s canceled the request", message.User.Name)
		responseMessage(w, message.OriginalMessage, title, "")
		return
	default:
		log.Printf("[ERROR] ]Invalid action was submitted: %s", action.Name)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// responseMessage response to the original slackbutton enabled message.
// It removes button and replace it with message which indicate how bot will work
func responseMessage(w http.ResponseWriter, original slack.Message, title, value string) {
	original.Attachments[0].Actions = []slack.AttachmentAction{} // empty buttons
	original.Attachments[0].Fields = []slack.AttachmentField{
		{
			Title: title,
			Value: value,
			Short: false,
		},
	}

	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&original)
}
