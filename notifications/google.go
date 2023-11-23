package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/chat/v1"
)

type GoogleChat struct {
	WebhookURL string
}

// NewGoogleChat returns a new instance of the GoogleChat struct.
//
// No parameters.
// Returns a pointer to a GoogleChat object.
func NewGoogleChat() *GoogleChat {
	return &GoogleChat{}
}

// SetWebhookURL sets the webhook URL for the GoogleChat instance.
//
// Parameters:
// - url: The URL to set as the webhook URL.
func (gc *GoogleChat) SetWebhookURL(url string) {
	gc.WebhookURL = url
}

// SendMessage sends a message to the Google Chat API webhook.
//
// It takes a string parameter `message` which is the content of the message to be sent.
// The function returns a boolean value indicating whether the message was sent successfully or not.
func (gc *GoogleChat) SendMessage(message string) bool {
	msg := chat.Message{
		Text: message,
	}
	// Convert message to JSON
	msgJson, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}
	// Send HTTP POST request to Google Chat webhook
	resp, err := http.Post(gc.WebhookURL, "application/json", bytes.NewBuffer(msgJson))
	if err != nil {
		log.Fatalf("Error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Message sent successfully!")
		return true
	} else {
		log.Fatalf("Failed to send message. Status code: %d\n", resp.StatusCode)
	}
	return false
}

// SendGridCard sends a Google Apps Card message to a Google Chat webhook.
//
// The function takes a title string and a gridMsgs [][]string as parameters.
// The title is used as the title of the card, and gridMsgs is a 2D slice
// containing the messages to be displayed in the card's grid format.
//
// The function returns a boolean value indicating whether the message was
// successfully sent or not.
func (gc *GoogleChat) SendGridCard(title string, gridMsgs [][]string) bool {
	card := chat.GoogleAppsCardV1Card{
		Header: &chat.GoogleAppsCardV1CardHeader{
			Title: title,
		},
		Sections: []*chat.GoogleAppsCardV1Section{
			{
				Collapsible: false,
				Widgets:     []*chat.GoogleAppsCardV1Widget{},
			},
		},
	}

	// Convert gridMsgs to card sections
	for _, row := range gridMsgs {
		r := chat.GoogleAppsCardV1Widget{
			Columns: &chat.GoogleAppsCardV1Columns{
				ColumnItems: []*chat.GoogleAppsCardV1Column{},
			},
		}
		for _, cell := range row {
			ci := chat.GoogleAppsCardV1Column{
				HorizontalSizeStyle: "FILL_AVAILABLE_SPACE",
				Widgets: []*chat.GoogleAppsCardV1Widgets{
					{
						TextParagraph: &chat.GoogleAppsCardV1TextParagraph{
							Text: cell,
						},
					},
				},
			}
			r.Columns.ColumnItems = append(r.Columns.ColumnItems, &ci)
		}

		card.Sections[0].Widgets = append(card.Sections[0].Widgets, &r)
	}

	msg := chat.Message{
		CardsV2: []*chat.CardWithId{
			{
				Card: &card,
			},
		},
	}

	// Convert message to JSON
	msgJson, err := json.Marshal(msg)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}

	// Send HTTP POST request to Google Chat webhook
	resp, err := http.Post(gc.WebhookURL, "application/json", bytes.NewBuffer(msgJson))
	if err != nil {
		log.Fatalf("Error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Message sent successfully!")
		return true
	} else {
		log.Fatalf("Failed to send message. Status code: %d\n", resp.StatusCode)
	}

	return false
}
