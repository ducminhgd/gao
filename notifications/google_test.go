package notifications

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"google.golang.org/api/chat/v1"
)

// TestNewGoogleChat is a test function that tests the functionality of the NewGoogleChat function.
//
// The function takes no parameters and returns nothing.
func TestNewGoogleChat(t *testing.T) {
	gc := NewGoogleChat()

	if gc == nil {
		t.Error("NewGoogleChat should not return nil")
	}
}

// TestSetWebhookURL is a test function that tests the SetWebhookURL method of the GoogleChat struct.
//
// It sets the webhook URL of the GoogleChat instance to the given URL and asserts that the value
// is correctly set.
// Parameter(s):
// - t: A *testing.T instance used for reporting test failures.
// Return type(s): None.
func TestSetWebhookURL(t *testing.T) {
	gc := &GoogleChat{}
	url := "https://example.com/webhook"

	gc.SetWebhookURL(url)

	if gc.WebhookURL != url {
		t.Errorf("expected %v, got %v", url, gc.WebhookURL)
	}
}

// TestSendMessage tests the SendMessage function of the GoogleChat struct.
//
// It creates a new http test server and sends a message to the server using
// the SendMessage method. The function checks that the HTTP method is POST,
// the content type is application/json, and the message sent is correct. It
// also checks that the HTTP response status is 200 and returns a boolean
// indicating whether the message was sent successfully.
func TestSendMessage(t *testing.T) {
	mockMessage := "Hello, World!"
	gc := &GoogleChat{}

	// Create a new http test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check that the method is POST
		if r.Method != http.MethodPost {
			t.Fatalf("Expected method 'POST', got '%s'", r.Method)
		}
		// Check that the content type is application/json
		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Expected content type 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}

		// Read the body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Body.Close()

		var msg chat.Message
		// Unmarshal the body into a chat.Message object
		json.Unmarshal(body, &msg)

		// Check that the message is correct
		if msg.Text != mockMessage {
			t.Fatalf("Expected message '%s', got '%s'", mockMessage, msg.Text)
		}

		// Respond with a status of 200
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	// Set the mock URL
	gc.SetWebhookURL(ts.URL)

	// Send the message
	success := gc.SendMessage(mockMessage)

	// Check that the message was sent successfully
	if !success {
		t.Errorf("Message was not sent successfully.")
	}
}

func TestSendGridCard(t *testing.T) {
	gc := &GoogleChat{}
	title := "Test Title"
	gridMsgs := [][]string{
		{"Row 1, Col 1", "Row 1, Col 2"},
		{"Row 2, Col 1", "Row 2, Col 2"},
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Fatalf("Expected method 'POST', got '%s'", r.Method)
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("Expected content type 'application/json', got '%s'", r.Header.Get("Content-Type"))
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		defer r.Body.Close()

		var msg chat.Message
		json.Unmarshal(body, &msg)

		if msg.CardsV2[0].Card.Header.Title != title {
			t.Fatalf("Expected title '%s', got '%s'", title, msg.CardsV2[0].Card.Header.Title)
		}

		if len(msg.CardsV2[0].Card.Sections[0].Widgets) != len(gridMsgs) {
			t.Fatalf("Expected %d widgets, got %d", len(gridMsgs), len(msg.CardsV2[0].Card.Sections[0].Widgets))
		}

		for i, widget := range msg.CardsV2[0].Card.Sections[0].Widgets {
			for j, item := range widget.Columns.ColumnItems {
				if item.Widgets[0].TextParagraph.Text != gridMsgs[i][j] {
					t.Fatalf("Expected grid cell '%s', got '%s'", gridMsgs[i][j], item.Widgets[0].TextParagraph.Text)
				}
			}
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	gc.SetWebhookURL(ts.URL)

	success := gc.SendGridCard(title, gridMsgs)

	if !success {
		t.Errorf("Grid card was not sent successfully.")
	}
}
