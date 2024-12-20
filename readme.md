### Simple Mistral Client for Golang

# Wrapper MistralApi


![Mistral Models Overview docs](https://docs.mistral.ai/getting-started/models/models_overview/)

### Example
# Sending message
```go
//init mistral client
client := mistal_client.New("API_KEY")

//set Message List
var mesageList []mistal_client.Message
mesageList = append(mesageList, mistal_client.Message{
Role:    "user",
Content: "посчитай до 5",
})
//send message
result, err := client.Mistral.SendMessage(
mistal_client.SendMessageRequest{
Model:    "mistral-large-latest",
Messages: mesageList,
})
if err != nil {
    fmt.Println(err)    
}

fmt.Println(result)
```


# Sending message stream

```go
//init mistral client
client := mistal_client.New("API_KEY")

//set Message List
var mesageList []mistal_client.Message
mesageList = append(mesageList, mistal_client.Message{
Role:    "user",
Content: "посчитай до 5",
})
//send message
iter, err := client.Mistral.SendMessageStream(ctx,
mistal_client.SendMessageRequest{
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
```