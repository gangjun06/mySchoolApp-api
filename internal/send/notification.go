package send

import (
	"fmt"

	expo "github.com/oliveroneill/exponent-server-sdk-golang/sdk"
)

func SendNotification(id, title, content string, data map[string]string) {
	pushToken, err := expo.NewExponentPushToken(id)
	if err != nil {
		panic(err)
	}

	client := expo.NewPushClient(nil)

	response, err := client.Publish(
		&expo.PushMessage{
			To:       []expo.ExponentPushToken{pushToken},
			Body:     content,
			Data:     data,
			Sound:    "default",
			Title:    title,
			Priority: expo.DefaultPriority,
		},
	)
	// Check errors
	if err != nil {
		fmt.Errorf("%s\n", err.Error())
		return
	}
	// Validate responses
	if response.ValidateResponse() != nil {
		fmt.Println(response.PushMessage.To, "failed")
	}
}
