package main

import (
	"context"
	"fmt"
	"github.com/Vic07Region/mistral_client"
	"time"
)

func main() {
	ai := mistral_client.New("API_KEY")

	var mesageList []mistral_client.Message

	mesageList = append(mesageList, mistral_client.Message{
		Role:    "user",
		Content: "посчитай до 5",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

	defer cancel()

	iter, err := ai.Mistral.SendMessageStream(ctx,
		mistral_client.SendMessageRequest{
			Model:    "mistral-large-latest",
			Messages: mesageList,
		},
	)
	if err != nil {
		fmt.Println(err)
	}

	for iter.Next() {
		if iter.Err() != nil {
			fmt.Println(iter.Err())
		}
		fmt.Printf("%v", iter.Value())
	}
}
