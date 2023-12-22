package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"strings"
)

func mentions(emails []string) []map[string]interface{} {
	m := []map[string]interface{}{}
	for _, email := range emails {
		m = append(m, map[string]interface{}{
			"type": "mention",
			"text": fmt.Sprintf("<at>%s</at>", email),
			"mentioned": map[string]interface{}{
				"id":   email,
				"name": email,
			},
		})
	}
	return m
}

func referenceEmails(emails []string) string {
	ref := []string{}
	for _, email := range emails {
		ref = append(ref, fmt.Sprintf("<at>%s</at>", email))
	}
	return strings.Join(ref, ",")
}

func messageBody(title, body string, notificationList []string) []map[string]interface{} {
	items := []map[string]interface{}{
		{
			"type":   "TextBlock",
			"size":   "Medium",
			"weight": "Bolder",
			"text":   title,
		},
		{
			"type": "TextBlock",
			"text": body,
		},
	}

	if len(notificationList) > 0 {
		items = append(items, map[string]interface{}{
			"type": "TextBlock",
			"text": fmt.Sprintf("cc %s", referenceEmails(notificationList)),
		})
	}
	return items
}

func createMessage(title, body string, notificationList []string) map[string]interface{} {
	return map[string]interface{}{
		"type": "message",
		"attachments": []map[string]interface{}{
			{
				"contentType": "application/vnd.microsoft.card.adaptive",
				"content": map[string]interface{}{
					"type":    "AdaptiveCard",
					"body":    messageBody(title, body, notificationList),
					"$schema": "http://adaptivecards.io/schemas/adaptive-card.json",
					"version": "1.0.0",
					"msteams": map[string]interface{}{
						"width":    "Full",
						"entities": mentions(notificationList),
					},
				},
			},
		},
	}
}

func main() {
	var (
		uri         string
		title       string
		body        string
		mentionList string
	)

	flag.StringVar(&uri, "uri", "", "MS Teams Webhook URI (mandatory)")
	flag.StringVar(&title, "title", "", "The message title")
	flag.StringVar(&body, "body", "", "The message body")
	flag.StringVar(&mentionList, "mentions", "", "A comma-separated link of mentions (in the form of emails) to be included in the message")

	flag.Parse()

	if uri == "" {
		log.Fatal("No webhook URL specified. Provide the -h flag to see program options.")
	}
	if title == "" {
		log.Fatal("No message title specified. Provide the -h flag to see program options.")
	}
	if body == "" {
		log.Fatal("No message body specified. Provide the -h flag to see program options.")
	}

	mentions := []string{}
	if mentionList != "" {
		mentions = strings.Split(mentionList, ",")
	}

	message := createMessage(title, body, mentions)

	payload, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("failed to serialize the message: %v", err)
	}

	resp, err := http.Post(uri, "application/json", bytes.NewReader(payload))
	if err != nil {
		log.Fatalf("failed to send the http POST request: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		log.Fatalf("service replied with status code: %d", resp.StatusCode)
	}

	log.Printf("Message posted successfully")
}
