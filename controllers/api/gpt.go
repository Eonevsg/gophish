package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/joho/godotenv"
)

type Question struct {
	HintText string `json:"hintText"`
}

// Gpt handles the functionality for the /api/gpt endpoint
func (as *Server) Gpt(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var question Question

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err)
		}
		json.Unmarshal(b, &question)

		response := GetResponseWithEngine(question.HintText)
		byteArray, err := json.Marshal(response)
		if err != nil {
			fmt.Println(err)
		}

		w.Write(byteArray)

	default:
		byteArray, err := json.Marshal("Sorry, only POST method is supported.")
		if err != nil {
			fmt.Println(err)
		}
		w.Write(byteArray)
	}
}

func GetResponseWithEngine(question string) string {
	godotenv.Load()

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatalln("Missing API Key")
	}

	ctx := context.Background()
	client := gpt3.NewClient(apiKey)

	var response string
	err := client.CompletionStreamWithEngine(ctx, gpt3.TextDavinci003Engine, gpt3.CompletionRequest{
		Prompt:      []string{question},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, func(resp *gpt3.CompletionResponse) {
		fmt.Print(resp.Choices[0].Text)
		response += resp.Choices[0].Text
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(13)
	}
	fmt.Printf("\n")
	return response
}
