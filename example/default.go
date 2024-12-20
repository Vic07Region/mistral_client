package main

import (
	"fmt"
	"github.com/Vic07Region/mistral_client"
)

func main() {
	ai := mistral_client.New("API_KEY")
	var mesageList []mistral_client.Message
	mesageList = append(mesageList, mistral_client.Message{
		Role:    "user",
		Content: "посчитай до 5",
	})

	result, err := ai.Mistral.SendMessage(
		mistral_client.SendMessageRequest{
			Model:    "pkg-large-latest",
			Messages: mesageList,
		},
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result)
}
