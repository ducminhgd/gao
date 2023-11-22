package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/chat/v1"
)

const webhookURL = "https://chat.googleapis.com/v1/spaces/xi0pzEAAAAE/messages?key=AIzaSyDdI0hCZtE6vySjMm-WEfRq3CPzqKqqsHI&token=S-yI8l7E4PPEAiq9Wgs66NB8fq_hPc2ltk6Gnzt_A6o"

type Message struct {
	Text string `json:"text"`
}

func main() {

	message := chat.Message{
		CardsV2: []*chat.CardWithId{
			{
				Card: &chat.GoogleAppsCardV1Card{
					Header: &chat.GoogleAppsCardV1CardHeader{
						Title: "Card Header",
					},
					Sections: []*chat.GoogleAppsCardV1Section{
						{
							Collapsible: false,
							Widgets: []*chat.GoogleAppsCardV1Widget{
								// {
								// 	TextParagraph: &chat.GoogleAppsCardV1TextParagraph{
								// 		Text: "*Card Text Paragraph*",
								// 	},
								// },
								// {
								// 	DecoratedText: &chat.GoogleAppsCardV1DecoratedText{
								// 		Text: "Content with <b>bold</b>, <i>italic</i>, and <s>strikethrough</s> text.",
								// 	},
								// },
								{
									Columns: &chat.GoogleAppsCardV1Columns{
										ColumnItems: []*chat.GoogleAppsCardV1Column{
											{
												Widgets: []*chat.GoogleAppsCardV1Widgets{
													{
														TextParagraph: &chat.GoogleAppsCardV1TextParagraph{
															Text: "<b>Key 1</b>",
														},
													},
												},
											},
											{
												HorizontalSizeStyle: "FILL_AVAILABLE_SPACE",
												Widgets: []*chat.GoogleAppsCardV1Widgets{
													{
														TextParagraph: &chat.GoogleAppsCardV1TextParagraph{
															Text: "Value 1",
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Convert message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Fatalf("Error encoding JSON: %v", err)
	}

	// Send HTTP POST request to Google Chat webhook
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(messageJSON))
	if err != nil {
		log.Fatalf("Error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Message sent successfully!")
	} else {
		fmt.Printf("Failed to send message. Status code: %d\n", resp.StatusCode)
	}
}
